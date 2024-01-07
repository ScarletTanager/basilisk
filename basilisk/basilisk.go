package main

import (
	"github.com/ScarletTanager/basilisk/basilisk/handlers"
	"github.com/ScarletTanager/basilisk/basilisk/model"
	"github.com/ScarletTanager/basilisk/classifiers/knn"
	"github.com/labstack/echo/v4"
)

var currentModel *knn.KNearestNeighborClassifier

func main() {
	rm := &model.RunningModels{}
	e := echo.New()

	e.POST("/models", handlers.CreateModelHandler(rm), handlers.CheckContentTypeMiddleware(handlers.AllowedHeaders{echo.MIMEApplicationJSON}))
	e.GET("/models", handlers.ListModelsHandler(rm))
	e.POST("/datasets", handlers.CreateDatasetHandler, handlers.CheckContentTypeMiddleware(handlers.AllowedHeaders{echo.MIMEApplicationJSON}))

	modelGroup := e.Group("/models/:id", handlers.RetrieveModelMiddleware(rm))
	modelGroup.PUT("/data", handlers.TrainModelHandler(rm), handlers.CheckContentTypeMiddleware(handlers.AllowedHeaders{echo.MIMEApplicationJSON, "text/csv"}))
	modelGroup.GET("/results", handlers.TestModelHandler(rm))
	modelGroup.GET("/results/details", handlers.TestResultsDetailsHandler(rm))

	e.Logger.Fatal(e.Start(":9323"))
}
