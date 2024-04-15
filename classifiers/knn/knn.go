package knn

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/sphinx/probability"
)

type KNearestNeighborClassifier struct {
	RawData       *classifiers.DataSet
	TrainingData  *classifiers.DataSet
	TestingData   *classifiers.DataSet
	Results       classifiers.TestResults
	Configuration KNearestNeighborClassifierConfig
}

type KNearestNeighborClassifierConfig struct {
	K                int
	DistanceMethod   string
	distanceFunction classifiers.DistanceFunction
}

func (knnc *KNearestNeighborClassifier) Config() interface{} {
	return knnc.Configuration
}

const (
	ClassifierType_KNearestNeighbor string = "KNearestNeighbors Classifier"
)

func (knnc *KNearestNeighborClassifier) Type() string {
	return ClassifierType_KNearestNeighbor
}

func New(k int, distanceMethod string) (*KNearestNeighborClassifier, error) {
	if k <= 0 {
		return nil, errors.New("Unable to create classifier, k must be greater than 0")
	}

	var distanceFunc classifiers.DistanceFunction

	switch distanceMethod {
	case classifiers.DistanceMethod_Euclidean:
		distanceFunc = classifiers.EuclideanDistance
	case classifiers.DistanceMethod_Manhattan:
		distanceFunc = classifiers.ManhattanDistance
	default:
		distanceMethod = classifiers.DistanceMethod_Euclidean
		distanceFunc = classifiers.EuclideanDistance
	}

	return &KNearestNeighborClassifier{
		Configuration: KNearestNeighborClassifierConfig{K: k, DistanceMethod: distanceMethod, distanceFunction: distanceFunc},
	}, nil
}

func (knnc *KNearestNeighborClassifier) Data() (*classifiers.DataSet, *classifiers.DataSet) {
	return knnc.TrainingData, knnc.TestingData
}

func (knnc *KNearestNeighborClassifier) TrainFromCSV(data []byte, cfg *classifiers.DataSplitConfig) error {
	return nil
}

func (knnc *KNearestNeighborClassifier) TrainFromCSVFile(path string, cfg *classifiers.DataSplitConfig) error {
	var err error
	knnc.RawData, err = classifiers.FromCSVFile(path)
	if err != nil {
		return fmt.Errorf("Error training from CSV file %s: %w", path, err)
	}

	return knnc.train(cfg)
}

func (knnc *KNearestNeighborClassifier) TrainFromDataset(ds *classifiers.DataSet, cfg *classifiers.DataSplitConfig) error {
	knnc.RawData = ds
	return knnc.train(cfg)
}

func (knnc *KNearestNeighborClassifier) TrainFromJSON(data []byte, cfg *classifiers.DataSplitConfig) error {
	return nil
}

func (knnc *KNearestNeighborClassifier) TrainFromJSONFile(path string, cfg *classifiers.DataSplitConfig) error {
	var err error
	knnc.RawData, err = classifiers.FromJSONFile(path)
	if err != nil {
		return fmt.Errorf("Error training from JSON file: %w", err)
	}

	return knnc.train(cfg)
}

func (knnc *KNearestNeighborClassifier) train(cfg *classifiers.DataSplitConfig) error {
	var err error
	knnc.TrainingData, knnc.TestingData, err = knnc.RawData.Split(cfg)
	return err
}

func (knnc *KNearestNeighborClassifier) Retrain(cfg *classifiers.DataSplitConfig) error {
	return knnc.train(cfg)
}

func (knnc *KNearestNeighborClassifier) Test() (classifiers.TestResults, error) {
	if knnc.TrainingData == nil || knnc.TestingData == nil {
		return nil, errors.New("Untestable model")
	}
	results := make(classifiers.TestResults, len(knnc.TestingData.Records))
	for i, testRecord := range knnc.TestingData.Records {
		results[i] = classify(testRecord,
			computeNeighbors(testRecord, knnc.TrainingData.Records, knnc.Configuration.distanceFunction),
			knnc.Configuration.K,
			len(knnc.RawData.ClassNames))
	}

	knnc.Results = results
	return results, nil
}

type Neighbor struct {
	Class    int
	Distance float64
}

// classify assumes that neighbors has been sorted by distance already
func classify(orig classifiers.Record, neighbors []Neighbor, k, classCount int) classifiers.TestResult {
	result := classifiers.TestResult{
		Record: orig,
	}

	votingNeighbors := neighbors[:k]
	votes := make([]int, k)

	// Collect the votes
	for ni, neighbor := range votingNeighbors {
		votes[ni] = neighbor.Class
	}

	// Create the probability mass function
	pmf := probability.MassDiscrete(votes)

	predicted := classifiers.NO_PREDICTION
	predictedProbability := 0.0

	// Determine the class
	for i := 0; i < classCount; i++ {
		probability := pmf(i)
		if probability > predictedProbability {
			predicted = i
			predictedProbability = probability
		}
	}

	result.Predicted = predicted
	result.Probability = predictedProbability

	return result
}

// Take an individual record, order the records from ds by proximity, return the ordered list
func computeNeighbors(orig classifiers.Record, comps []classifiers.Record, distanceFunction classifiers.DistanceFunction) []Neighbor {
	neighbors := make([]Neighbor, len(comps))
	for i, r := range comps {
		neighbors[i] = Neighbor{
			Class:    r.Class,
			Distance: distanceFunction(orig.AttributeValues, r.AttributeValues),
		}
	}

	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].Distance < neighbors[j].Distance
	})

	return neighbors
}
