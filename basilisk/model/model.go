package model

import (
	"errors"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/basilisk/classifiers/knn"
)

type ModelConfiguration struct {
	K int `json:"k,omitempty"`
}

type Model struct {
	ID int `json:"id"`
	ModelConfiguration
}

type RunningModels struct {
	Classifiers []classifiers.Classifier
}

type ModelsError struct {
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}

func (rm *RunningModels) Add(cl *knn.KNearestNeighborClassifier) (int, error) {
	if cl == nil {
		return -1, errors.New("Cannot add a nil classifier")
	}

	if rm.Classifiers == nil {
		rm.Classifiers = make([]classifiers.Classifier, 0)
	}

	rm.Classifiers = append(rm.Classifiers, cl)
	return len(rm.Classifiers) - 1, nil
}
