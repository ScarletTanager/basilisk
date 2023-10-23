package handlers

import (
	"errors"
	"net/http"

	"github.com/ScarletTanager/basilisk/classifiers/knn"
	"github.com/labstack/echo/v4"
)

type ModelConfiguration struct {
	K int `json:"k,omitempty"`
}

type Model struct {
	ID int `json:"id"`
	ModelConfiguration
}

type RunningModels struct {
	Classifiers []*knn.KNearestNeighborClassifier
}

type ModelsError struct {
	Message string
}

func (rm *RunningModels) Add(cl *knn.KNearestNeighborClassifier) (int, error) {
	if cl == nil {
		return -1, errors.New("Cannot add a nil classifier")
	}

	if rm.Classifiers == nil {
		rm.Classifiers = make([]*knn.KNearestNeighborClassifier, 0)
	}

	rm.Classifiers = append(rm.Classifiers, cl)
	return len(rm.Classifiers) - 1, nil
}

// CreateModelHandler returns an echo.HandlerFunc configured to set the currentModel with a valid request
func CreateModelHandler(rm *RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id  int
			err error
		)

		mc := new(ModelConfiguration)
		if err = c.Bind(mc); err != nil {
			return c.JSON(http.StatusBadRequest, &ModelsError{Message: "Cannot parse request body"})
		}

		classifier, _ := knn.New(mc.K)
		if id, err = rm.Add(classifier); err != nil {
			return c.JSON(http.StatusInternalServerError, &ModelsError{Message: "Server error creating model, please retry"})
		}

		return c.JSON(http.StatusOK, &Model{ID: id, ModelConfiguration: ModelConfiguration{mc.K}})
	}
}
