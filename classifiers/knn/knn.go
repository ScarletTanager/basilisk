package knn

import (
	"fmt"

	"github.com/ScarletTanager/basilisk/classifiers"
)

type KNearestNeighborClassifier struct {
	RawData      *classifiers.DataSet
	TrainingData *classifiers.DataSet
	TestingData  *classifiers.DataSet
}

func New() *KNearestNeighborClassifier {
	return &KNearestNeighborClassifier{}
}

func (knnc *KNearestNeighborClassifier) TrainFromJSONFile(path string) error {
	var err error
	knnc.RawData, err = classifiers.FromJSONFile(path)
	if err != nil {
		return fmt.Errorf("Error training from JSON file: %w", err)
	}

	return nil
}

func (knnc *KNearestNeighborClassifier) Retrain() {

}

func (knnc *KNearestNeighborClassifier) Test() classifiers.TestResults {
	return nil
}
