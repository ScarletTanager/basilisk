package knn_test

import (
	"github.com/ScarletTanager/basilisk/classifiers/knn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Knn", func() {
	var (
		knnc *knn.KNearestNeighborClassifier
	)

	BeforeEach(func() {
		knnc = knn.New()
	})

	Describe("TrainFromJSONFile", func() {
		var (
			path string
		)

		BeforeEach(func() {
			path = "../../datasets/shorebirds.json"
		})

		It("Does not return an error", func() {
			Expect(knnc.TrainFromJSONFile(path)).NotTo(HaveOccurred())
		})

		When("It cannot read the data file", func() {
			BeforeEach(func() {
				path = "../../datasets/thisdoesnotexist.json"
			})

			It("Returns an error", func() {
				Expect(knnc.TrainFromJSONFile(path)).To(HaveOccurred())
			})
		})
	})
})
