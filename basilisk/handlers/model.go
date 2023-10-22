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

type RunningModel struct {
	Classifier *knn.KNearestNeighborClassifier
}

func (rm *RunningModel) Set(cl *knn.KNearestNeighborClassifier) error {
	if rm.Classifier != nil {
		return errors.New("current classifier already set")
	}

	rm.Classifier = cl
	return nil
}

// CreateModelHandler returns an echo.HandlerFunc configured to set the currentModel with a valid request
func CreateModelHandler(rm *RunningModel) echo.HandlerFunc {
	return func(c echo.Context) error {
		mc := new(ModelConfiguration)
		if err := c.Bind(mc); err != nil {
			return err
		}

		classifier, _ := knn.New(mc.K)
		if err := rm.Set(classifier); err != nil {
			return c.String(http.StatusBadRequest, "Model exists already")
		}

		return c.String(http.StatusOK, "Model created")
	}
}
