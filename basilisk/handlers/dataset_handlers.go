package handlers

import (
	"net/http"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/basilisk/dsgen"
	"github.com/labstack/echo/v4"
)

//
// Handlers for working with datasets
//

func CreateDatasetHandler(c echo.Context) error {
	var (
		dataset *classifiers.DataSet
		err     error
	)

	datasetConfig := &dsgen.DatasetConfig{}
	if err = c.Bind(datasetConfig); err != nil {
		return c.JSON(http.StatusBadRequest, &dsgen.DatasetConfigError{
			Message: "Unable to process configuration",
		})
	}

	if dataset, err = dsgen.GenerateDataset(datasetConfig); err != nil {
		return c.JSON(http.StatusBadRequest, &dsgen.DatasetConfigError{
			Err:     err,
			Message: "Unable to generate dataset due to errors in configuration",
		})
	}

	return c.JSON(http.StatusOK, dataset)
}
