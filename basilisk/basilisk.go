package main

import (
	"github.com/ScarletTanager/basilisk/basilisk/handlers"
	"github.com/ScarletTanager/basilisk/classifiers/knn"
	"github.com/labstack/echo/v4"
)

var currentModel *knn.KNearestNeighborClassifier

func main() {
	rm := &handlers.RunningModels{}
	e := echo.New()

	e.POST("/models", handlers.CreateModelHandler(rm))
	e.PUT("/models/:id/trainingdata", handlers.TrainModelHandler(rm))
	e.Logger.Fatal(e.Start(":9323"))
}
