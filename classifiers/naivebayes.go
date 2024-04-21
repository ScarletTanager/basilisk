package classifiers

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/ScarletTanager/sphinx/probability"
)

//
// Simple implementation of a naive Bayes classifier
//

type NaiveBayesClassifier struct {
	ClassifierImplementation
	VectorConditionedClassProbabilities [][]float64
}

const (
	ClassifierType_NaiveBayes string = "Naive Bayes Classifier"
)

func NewNaiveBayes() *NaiveBayesClassifier {
	return &NaiveBayesClassifier{}
}

func (nbc *NaiveBayesClassifier) Type() string {
	return ClassifierType_NaiveBayes
}

func (nbc *NaiveBayesClassifier) Config() interface{} {
	return struct{}{}
}

func (nbc *NaiveBayesClassifier) TrainFromCSVFile(path string, cfg *DataSplitConfig) error {
	var err error
	nbc.RawData, err = FromCSVFile(path)
	if err != nil {
		return fmt.Errorf("Error training from CSV file %s: %w", path, err)
	}

	return nbc.train(cfg)
}

func (nbc *NaiveBayesClassifier) TrainFromDataset(ds *DataSet, cfg *DataSplitConfig) error {
	nbc.RawData = ds
	return nbc.train(cfg)
}

func (nbc *NaiveBayesClassifier) train(cfg *DataSplitConfig) error {
	trainingData, testingData, err := nbc.RawData.Split(cfg)
	if err != nil {
		return err
	}

	nbc.TrainingData = trainingData
	nbc.TestingData = testingData

	// Calculate the prior probability of each class
	// This should probably be a method on the DataSet...?

	classPriors := ComputeClassPriors(nbc.RawData.ClassNames, nbc.TrainingData.Records)

	// Discretize the attributes - we will use the output from this step
	// to calculate the conditional probability P(X==x|C==c), where x is the vector of attributes,
	// and c is the class.

	attributeIntervals := DiscretizeAttributes(nbc.RawData.AttributeNames, nbc.TrainingData.Records)

	// Now we need to calcuate P(An==a|C==c) (A1...An are the individual attribute variables)

	// For each class, for each attribute, we compute the probability for
	// each interval

	caps := GenerateClassAttributePosteriors(nbc.RawData.ClassNames, nbc.RawData.AttributeNames,
		attributeIntervals, nbc.TrainingData.Records)

	// Next calculate P(X==x|C==c)

	ccps := GenerateClassConditionedPosteriors(nbc.RawData.ClassNames, nbc.RawData.AttributeNames,
		attributeIntervals, caps)

	// Aaaaand...apply Bayes to get P(C==c|X==x)

	vectorCount := len(ccps[0])

	fmt.Fprintf(os.Stderr, "Vector count: %d\n", vectorCount)

	nbc.VectorConditionedClassProbabilities = make([][]float64, vectorCount)
	for v, _ := range nbc.VectorConditionedClassProbabilities {
		nbc.VectorConditionedClassProbabilities[v] = make([]float64, len(nbc.RawData.ClassNames))
	}

	vectorTotalProbabilities := make([]float64, vectorCount)
	for i, _ := range vectorTotalProbabilities {
		for classIdx, _ := range nbc.RawData.ClassNames {
			vectorTotalProbabilities[i] += ccps[classIdx][i] * classPriors[classIdx]
		}

		fmt.Fprintf(os.Stderr, "Vector: %d\tTotal probability: %f\n", i, vectorTotalProbabilities[i])
	}

	for classIdx, _ := range ccps {
		for vectorIdx, pVal := range ccps[classIdx] {
			nbc.VectorConditionedClassProbabilities[vectorIdx][classIdx], _ = probability.Bayes(classPriors[classIdx],
				pVal, vectorTotalProbabilities[vectorIdx])
			fmt.Fprintf(os.Stderr, "Class: %d\tClass Prior: %.4f\tVector: %d\tVector Posterior: %.4f\tVector total: %.4f\tClass Posterior: %.4f\n",
				classIdx, classPriors[classIdx], vectorIdx, ccps[classIdx][vectorIdx], vectorTotalProbabilities[vectorIdx], nbc.VectorConditionedClassProbabilities[vectorIdx][classIdx])
			// fmt.Fprintf(os.Stderr, "Vector: %d\tClass: %d\tVector-conditioned Class posterior: %f\n", vectorIdx, classIdx, nbc.VectorConditionedClassProbabilities[vectorIdx][classIdx])
		}
	}

	return nil
}

func ComputeClassPriors(classNames []string, records []Record) []float64 {
	// Use the raw data in case the split left some classes out of the training dataset
	classPriors := make([]float64, len(classNames))

	// Stick the classes into a slice of ints
	trainingClasses := make([]int, len(records))
	for i, rec := range records {
		trainingClasses[i] = rec.Class
	}

	// Compute the class priors and cache them
	classPriorPMF := probability.MassDiscrete(trainingClasses)
	for i, _ := range classNames {
		classPriors[i] = classPriorPMF(i)
	}

	return classPriors
}

func DiscretizeAttributes(attributeNames []string, records []Record) []probability.Intervals {
	attributeIntervals := make([]probability.Intervals, len(attributeNames))

	// Range over the attributes
	for attrIdx, _ := range attributeNames {
		// Range over the records, store attribute-specific slices of values
		vals := make([]float64, len(records))
		for i, rec := range records {
			vals[i] = rec.AttributeValues[attrIdx]
		}

		//TODO: Make the number of intervals configurable when creating the classifier
		attributeIntervals[attrIdx] = probability.Discretize(vals, probability.DiscretizationConfig{})
	}

	return attributeIntervals
}

// ClassAttributePosteriors is just a convenience type alias to make retrieving the conditional
// probabilities easier
type ClassAttributePosteriors [][][]float64

func (cap ClassAttributePosteriors) ClassAttributePosterior(classIdx, attrIdx, intervalIdx int) (float64, error) {
	if cap[classIdx] == nil || cap[classIdx][attrIdx] == nil {
		return 0.0, errors.New("unable to retrieve probability, storage has not been initialized")
	}

	return cap[classIdx][attrIdx][intervalIdx], nil
}

// GenerateClassAttributePosteriors computes, for each class and attribute, the conditional probability
// that the attribute has a value in each interval, given the class.
func GenerateClassAttributePosteriors(classNames, attributeNames []string, intervals []probability.Intervals, records []Record) ClassAttributePosteriors {
	// Create the probability storage
	caps := make(ClassAttributePosteriors, len(classNames))

	for i, _ := range caps {
		caps[i] = make([][]float64, len(attributeNames))
		for j, _ := range caps[i] {
			caps[i][j] = make([]float64, len(intervals[j]))
		}
	}

	// Now convert the attribute values to intervals, first
	// create some storage for the interval values themselves
	cais := make([][][]int, len(classNames))
	for classIdx, _ := range cais {
		cais[classIdx] = make([][]int, len(attributeNames))
		for attrIdx, _ := range cais[classIdx] {
			cais[classIdx][attrIdx] = make([]int, 0)
		}
	}

	// Do the conversion and store the values
	for _, rec := range records {
		for attrIdx, _ := range attributeNames {
			cais[rec.Class][attrIdx] = append(cais[rec.Class][attrIdx],
				intervals[attrIdx].IntervalForValue(rec.AttributeValues[attrIdx]))
		}
	}

	// Create the PMFs
	classAttrPMFs := make([][]probability.ProbabilityMassFunction, len(classNames))
	for classIdx, _ := range classAttrPMFs {
		classAttrPMFs[classIdx] = make([]probability.ProbabilityMassFunction, len(attributeNames))
		for attrIdx, _ := range classAttrPMFs[classIdx] {
			classAttrPMFs[classIdx][attrIdx] = probability.MassDiscrete(cais[classIdx][attrIdx])
		}
	}

	// Compute the class-conditioned posteriors for each individual attribute
	for classIdx, _ := range caps {
		for attrIdx, _ := range caps[classIdx] {
			for intervalIdx, _ := range caps[classIdx][attrIdx] {
				caps[classIdx][attrIdx][intervalIdx] =
					classAttrPMFs[classIdx][attrIdx](intervalIdx)
			}
		}
	}

	return caps
}

type ClassConditionedPosteriors [][]float64

func GenerateClassConditionedPosteriors(classNames, attributeNames []string,
	intervals []probability.Intervals, caps ClassAttributePosteriors) ClassConditionedPosteriors {
	// We use 10 intervals per attribute, if we make that configurable,
	// then the base of the exponential must be the interval count
	vectorCount := int(math.Pow10(len(attributeNames)))

	ccps := make(ClassConditionedPosteriors, len(classNames))
	for classIdx, _ := range ccps {
		ccps[classIdx] = make([]float64, vectorCount)
		for i, _ := range ccps[classIdx] {
			ccps[classIdx][i] = 1.0
		}

		for attrIdx, _ := range attributeNames {
			chunkSize := int(math.Pow10(len(attributeNames) - 1 - attrIdx))
			vectorIdx := 0
			for vectorIdx < vectorCount {
				// If we ever allow the interval count to be configurable, this will need to change
				for k, _ := range intervals[attrIdx] {
					pVal, _ := caps.ClassAttributePosterior(classIdx, attrIdx, k)
					// Multiply the probability of every vector in the chunk by the interval posterior
					for chunkSlotIdx := 0; chunkSlotIdx < chunkSize; chunkSlotIdx++ {
						ccps[classIdx][vectorIdx] = ccps[classIdx][vectorIdx] * pVal
						vectorIdx++
					}
				}
			}
		}
	}

	return ccps
}
