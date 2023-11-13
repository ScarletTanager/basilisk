package handlers

import (
	"fmt"
	"net/http"

	"github.com/ScarletTanager/basilisk/basilisk/model"
	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/basilisk/classifiers/knn"
	"github.com/labstack/echo/v4"
)

// CreateModelHandler returns an echo.HandlerFunc configured to set the currentModel with a valid request
func CreateModelHandler(rm *model.RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id  int
			err error
		)

		mc := new(model.ModelConfiguration)
		if err = c.Bind(mc); err != nil {
			return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: "Cannot parse request body"})
		}

		classifier, _ := knn.New(mc.K)
		if id, err = rm.Add(classifier); err != nil {
			return c.JSON(http.StatusInternalServerError, &model.ModelsError{Message: "Server error creating model, please retry"})
		}

		return c.JSON(http.StatusOK, &model.Model{ID: id, ModelConfiguration: *mc})
	}
}

// TrainModelHandler returns an echo.HandlerFunc configured to use the uploaded data to train the specified
// model
func TrainModelHandler(rm *model.RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			err error
			raw classifiers.DataSet
			ds  *classifiers.DataSet
		)

		if knnc, ok := c.Get(ContextKeyModel).(classifiers.Classifier); !ok {
			return c.JSON(http.StatusNotFound, &model.ModelsError{Message: "Model not found"})
		} else {
			if err = c.Bind(&raw); err != nil {
				return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: "Cannot parse body"})
			}

			if ds, err = classifiers.NewDataSet(raw.ClassNames, raw.AttributeNames, raw.Records); err != nil {
				return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: fmt.Sprintf("Invalid data: %s", err.Error())})
			}

			if err = knnc.TrainFromDataset(ds, nil); err != nil {
				return c.JSON(http.StatusInternalServerError, &model.ModelsError{Message: "Error training model, please retry"})
			}
		}

		return c.JSON(http.StatusOK, []byte(`{"Message": "Completed"}`))
	}
}

func TestModelHandler(rm *model.RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
