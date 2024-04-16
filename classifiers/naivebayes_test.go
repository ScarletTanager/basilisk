package classifiers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/classifiers"
)

var _ = Describe("Naivebayes", func() {
	var (
		nbc  *classifiers.NaiveBayesClassifier
		cfg  *classifiers.DataSplitConfig
		path string
	)

	BeforeEach(func() {
		nbc = classifiers.NewNaiveBayes()
		cfg = nil
	})

	FContext("TrainFromCSVFile", func() {
		BeforeEach(func() {
			path = "../datasets/shorebirds.csv"
		})

		It("Does not return an error", func() {
			Expect(nbc.TrainFromCSVFile(path, cfg)).NotTo(HaveOccurred())
		})

		When("It cannot read the data file", func() {
			BeforeEach(func() {
				path = "../datasets/thisdoesnotexist.csv"
			})

			It("Returns an error", func() {
				Expect(nbc.TrainFromCSVFile(path, cfg)).To(HaveOccurred())
			})
		})
	})
})
