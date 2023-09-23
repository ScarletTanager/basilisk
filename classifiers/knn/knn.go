package knn

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ScarletTanager/basilisk/classifiers"
)

type KNearestNeighborClassifier struct {
	NeighborCount int
	RawData       *classifiers.DataSet
	TrainingData  *classifiers.DataSet
	TestingData   *classifiers.DataSet
	K             int
}

func New(k int) (*KNearestNeighborClassifier, error) {
	if k == 0 {
		return nil, errors.New("Unable to create classifier, k must be greater than 0")
	}
	return &KNearestNeighborClassifier{K: k}, nil
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

func (knnc *KNearestNeighborClassifier) Test() classifiers.TestResults {
	results := make(classifiers.TestResults, len(knnc.TestingData.Records))
	for i, testRecord := range knnc.TestingData.Records {
		results[i] = classify(testRecord, computeNeighbors(testRecord, knnc.TrainingData.Records), len(knnc.TestingData.Records)/2)
	}
	return results
}

type Neighbor struct {
	Class    int
	Distance float64
}

// classify assumes that neighbors has been sorted by distance already
func classify(orig classifiers.Record, neighbors []Neighbor, k int) classifiers.TestResult {
	var (
		winningVoteCount int
	)

	result := classifiers.TestResult{
		Record:    orig,
		Predicted: classifiers.NO_PREDICTION,
	}

	votingNeighbors := neighbors[:k]
	// Using a map because we don't know how many classes are represented in the set overall,
	// nor do we know which of those are represented in the training data
	votes := make(map[int]int)

	for _, n := range votingNeighbors {
		if _, ok := votes[n.Class]; !ok {
			votes[n.Class] = 1
		} else {
			votes[n.Class] = votes[n.Class] + 1
		}
	}

	for class, voteCount := range votes {
		if voteCount > winningVoteCount {
			result.Predicted = class
		}
	}

	return result
}

// Take an individual record, order the records from ds by proximity, return the ordered list
func computeNeighbors(orig classifiers.Record, comps []classifiers.Record) []Neighbor {
	neighbors := make([]Neighbor, len(comps))
	for i, r := range comps {
		neighbors[i] = Neighbor{
			Class:    r.Class,
			Distance: orig.AttributeValues.Difference(r.AttributeValues).Magnitude(),
		}
	}

	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].Distance < neighbors[j].Distance
	})

	return neighbors
}
