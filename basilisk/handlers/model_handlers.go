package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ScarletTanager/basilisk/basilisk/model"
	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ModelRenderer struct {
	ID                  int         `json:"id"`
	Type                string      `json:"type"`
	Config              interface{} `json:"config,omitempty"`
	Classes             []string    `json:"classes,omitempty"`
	Attributes          []string    `json:"attributes,omitempty"`
	TrainingRecordCount int         `json:"training_dataset_size"`
	TestingRecordCount  int         `json:"testing_dataset_size"`
}

func ListModelsHandler(rm *model.RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		if rm == nil {
			return c.JSON(http.StatusServiceUnavailable, []byte("Server not ready"))
		}

		body := make([]ModelRenderer, 0)
		for id, m := range rm.Classifiers {
			mr := ModelRenderer{
				ID:     id,
				Type:   m.Type(),
				Config: m.Config(),
			}

			trd, ted := m.Data()
			if trd != nil {
				mr.Classes = trd.Classes()
				mr.Attributes = trd.Attributes()
				mr.TrainingRecordCount = len(trd.Records)
			}

			if ted != nil {
				mr.TestingRecordCount = len(ted.Records)
			}

			body = append(body, mr)
		}

		return c.JSON(http.StatusOK, body)
	}
}

// CreateModelHandler returns an echo.HandlerFunc configured to set the currentModel with a valid request
func CreateModelHandler(rm *model.RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			createdModel *model.Model
		)

		mc := new(model.ModelConfiguration)
		if err := c.Bind(mc); err != nil {
			log.Errorf("Body binding error: %s", err.Error())
			return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: "Cannot parse request body", Error: err})
		}

		log.Infof("Model configuration: %v", mc)

		if classifier, err := classifiers.NewKnn(mc.K, mc.DistanceMethod); err != nil {
			return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: "Invalid model configuration", Error: err})
		} else {
			if id, err := rm.Add(classifier); err != nil {
				log.Errorf("Model creation error: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, &model.ModelsError{Message: "Server error creating model, please retry", Error: err})
			} else {
				createdModel = &model.Model{ID: id, ModelConfiguration: *mc}
			}
		}

		return c.JSON(http.StatusOK, createdModel)
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
			// Check what data format is being sent
			switch c.Request().Header.Get(echo.HeaderContentType) {
			case echo.MIMEApplicationJSON:
				if err = c.Bind(&raw); err != nil {
					return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: "Cannot parse body"})
				}

				if ds, err = classifiers.NewDataSet(raw.ClassNames, raw.AttributeNames, raw.Records); err != nil {
					return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: fmt.Sprintf("Invalid data: %s", err.Error())})
				}
			case "text/csv":
				bodyBytes, err := io.ReadAll(c.Request().Body)
				if err != nil {
					// Probably not really what we want here, but...whatever
					return c.JSON(http.StatusBadRequest, "Unable to read body")
				}
				if ds, err = classifiers.FromCSV(bodyBytes); err != nil || ds == nil {
					return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: fmt.Sprintf("Invalid data: %s", err.Error())})
				}
			}

			if err = knnc.TrainFromDataset(ds, nil); err != nil {
				return c.JSON(http.StatusInternalServerError, &model.ModelsError{Message: "Error training model, please retry"})
			}
		}

		return c.JSON(http.StatusOK, `{"Message": "Completed"}`)
	}
}

func TestModelHandler(rm *model.RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			tra classifiers.TestResultsAnalysis
		)

		if knnc, ok := c.Get(ContextKeyModel).(classifiers.Classifier); !ok {
			return c.JSON(http.StatusNotFound, &model.ModelsError{Message: "Model not found"})
		} else {
			if results, err := knnc.Test(); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			} else {
				tra = results.Analyze()
			}
		}
		return c.JSON(http.StatusOK, tra)
	}
}

func TestResultsDetailsHandler(rm *model.RunningModels) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			results classifiers.TestResults
			err     error
		)

		if knnc, ok := c.Get(ContextKeyModel).(classifiers.Classifier); !ok {
			return c.JSON(http.StatusNotFound, &model.ModelsError{Message: "Model not found"})
		} else {
			if results, err = knnc.Test(); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
		}

		return c.JSON(http.StatusOK, results)
	}

}
