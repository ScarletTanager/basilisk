package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ScarletTanager/basilisk/classifiers"
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

// TrainModelHandler returns an echo.HandlerFunc configured to use the uploaded data to train the specified
// model
func TrainModelHandler(rm *RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id  int
			err error
			raw classifiers.DataSet
			ds  *classifiers.DataSet
		)

		if id, err = strconv.Atoi(c.Param("id")); err != nil {
			return c.JSON(http.StatusBadRequest, &ModelsError{Message: fmt.Sprintf("%s is not a valid model id", c.Param("id"))})
		}

		if id < 0 || id > len(rm.Classifiers)-1 {
			return c.JSON(http.StatusNotFound, &ModelsError{Message: "Model not found"})
		}

		if err = c.Bind(&raw); err != nil {
			return c.JSON(http.StatusBadRequest, &ModelsError{Message: "Cannot parse body"})
		}

		if ds, err = classifiers.NewDataSet(raw.ClassNames, raw.AttributeNames, raw.Records); err != nil {
			return c.JSON(http.StatusBadRequest, &ModelsError{Message: fmt.Sprintf("Invalid data: %s", err.Error())})
		}

		rm.Classifiers[id].RawData = ds
		if err = rm.Classifiers[id].Retrain(nil); err != nil {
			return c.JSON(http.StatusInternalServerError, &ModelsError{Message: "Error training model, please retry"})
		}

		return c.JSON(http.StatusOK, &ModelsError{Message: "Model trained"})
	}
}
