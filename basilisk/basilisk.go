package main

import (
	"github.com/ScarletTanager/basilisk/basilisk/handlers"
	"github.com/ScarletTanager/basilisk/classifiers/knn"
	"github.com/labstack/echo"
)

var currentModel *knn.KNearestNeighborClassifier

func main() {
	e := echo.New()
	groupModels := e.Group("/models")
	groupModels.POST("/", handlers.CreateModelHandler(currentModel))
}
