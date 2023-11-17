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

	e.POST("/models", handlers.CreateModelHandler(rm), handlers.CheckContentTypeJSONMiddleware)
	e.GET("/models", handlers.ListModelsHandler(rm))
	e.POST("/datasets", handlers.CreateDatasetHandler, handlers.CheckContentTypeJSONMiddleware)

	modelGroup := e.Group("/models/:id", handlers.RetrieveModelMiddleware(rm))
	modelGroup.PUT("/data", handlers.TrainModelHandler(rm), handlers.CheckContentTypeJSONMiddleware)
	modelGroup.GET("/results", handlers.TestModelHandler(rm))

	e.Logger.Fatal(e.Start(":9323"))
}
