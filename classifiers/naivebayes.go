package classifiers

import (
	"fmt"

	"github.com/ScarletTanager/sphinx/probability"
)

//
// Simple implementation of a naive Bayes classifier
//

type NaiveBayesClassifier struct {
	ClassifierImplementation
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

	// Use the raw data in case the split left some classes out of the training dataset
	classPriors := make([]float64, len(nbc.RawData.ClassNames))

	// Stick the classes into a slice of ints
	trainingClasses := make([]int, len(nbc.TrainingData.Records))
	for i, rec := range nbc.TrainingData.Records {
		trainingClasses[i] = rec.Class
	}

	// Compute the class priors
	classPriorPMF := probability.MassDiscrete(trainingClasses)
	for i, _ := range nbc.RawData.ClassNames {
		classPriors[i] = classPriorPMF(i)
	}

	// Discretize the attributes - we will use the output from this step
	// to calculate the conditional probability P(X==x|C==c), where x is the vector of attributes,
	// and c is the class.

	attributeIntervals := make([]probability.Intervals, len(nbc.TrainingData.AttributeNames))

	// Range over the attributes
	for attrIdx, _ := range nbc.TrainingData.AttributeNames {
		// Range over the records to find min and max values
		vals := make([]float64, len(nbc.TrainingData.Records))
		for i, rec := range nbc.TrainingData.Records {
			vals[i] = rec.AttributeValues[attrIdx]
		}

		//TODO: Make the number of intervals configurable when creating the classifier
		attributeIntervals[attrIdx] = probability.Discretize(vals, probability.DiscretizationConfig{})
	}

	// Now we need to calcuate P(X==x|C==c)

	// For each class, for each attribute, we compute the probability for
	// each interval

	// Create the probability storage
	classAttributePosteriors := make([][][]float64, len(nbc.RawData.ClassNames))
	for i, _ := range classAttributePosteriors {
		classAttributePosteriors[i] = make([][]float64, len(nbc.RawData.AttributeNames))
		for j, _ := range classAttributePosteriors[i] {
			classAttributePosteriors[i][j] = make([]float64, len(attributeIntervals[j]))
		}
	}

	// Now convert the attribute values to intervals, first
	// create some storage for the interval values themselves
	classAttrIntervals := make([][][]int, len(nbc.RawData.ClassNames))
	for classIdx, _ := range nbc.RawData.ClassNames {
		classAttrIntervals[classIdx] = make([][]int, len(nbc.RawData.AttributeNames))
		for attrIdx, _ := range nbc.RawData.AttributeNames {
			classAttrIntervals[classIdx][attrIdx] = make([]int, 0)
		}
	}

	// Do the conversion and store the values
	for _, rec := range nbc.TrainingData.Records {
		for attrIdx, _ := range nbc.TrainingData.AttributeNames {
			classAttrIntervals[rec.Class][attrIdx] = append(classAttrIntervals[rec.Class][attrIdx],
				attributeIntervals[attrIdx].IntervalForValue(rec.AttributeValues[attrIdx]))
		}
	}

	// Create the PMFs
	classAttrPMFs := make([][]probability.ProbabilityMassFunction, len(nbc.RawData.ClassNames))
	for classIdx, _ := range classAttrPMFs {
		classAttrPMFs[classIdx] = make([]probability.ProbabilityMassFunction, len(nbc.RawData.AttributeNames))
		for attrIdx, _ := range classAttrPMFs[classIdx] {
			classAttrPMFs[classIdx][attrIdx] = probability.MassDiscrete(classAttrIntervals[classIdx][attrIdx])
		}
	}

	// Compute the class-conditioned posteriors for each individual attribute
	for classIdx, _ := range classAttributePosteriors {
		for attrIdx, _ := range classAttributePosteriors[classIdx] {
			for intervalIdx, _ := range classAttributePosteriors[classIdx][attrIdx] {
				classAttributePosteriors[classIdx][attrIdx][intervalIdx] =
					classAttrPMFs[classIdx][attrIdx](intervalIdx)
			}
		}
	}

	return nil
}
