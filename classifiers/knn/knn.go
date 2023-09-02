package knn

import "github.com/ScarletTanager/basilisk/classifiers"

type KNearestNeighborClassifier struct {
	TrainingData classifiers.DataSet
}

func New() *KNearestNeighborClassifier {
	return nil
}

// func (k *KNearestNeighborClassifier) Train()
