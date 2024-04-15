package classifiers_test

import (
	"github.com/ScarletTanager/basilisk/classifiers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Knn", func() {
	var (
		knnc           *classifiers.KNearestNeighborClassifier
		path           string
		cfg            *classifiers.DataSplitConfig
		k              int
		distanceMethod string
	)

	BeforeEach(func() {
		k = 3
		distanceMethod = classifiers.DistanceMethod_Euclidean
		cfg = nil
	})

	JustBeforeEach(func() {
		knnc, _ = classifiers.NewKnn(k, distanceMethod)
	})

	Describe("New", func() {
		It("Returns a new classifier", func() {
			c, err := classifiers.NewKnn(k, distanceMethod)
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(BeNil())
		})

		When("Called with k=0", func() {
			BeforeEach(func() {
				k = 0
			})

			It("Returns nil and an error", func() {
				c, err := classifiers.NewKnn(k, distanceMethod)
				Expect(err).To(HaveOccurred())
				Expect(c).To(BeNil())
			})
		})
	})

	Describe("TrainFromJSONFile", func() {
		BeforeEach(func() {
			path = "../datasets/shorebirds.json"
		})

		It("Does not return an error", func() {
			Expect(knnc.TrainFromJSONFile(path, cfg)).NotTo(HaveOccurred())
		})

		When("It cannot read the data file", func() {
			BeforeEach(func() {
				path = "../datasets/thisdoesnotexist.json"
			})

			It("Returns an error", func() {
				Expect(knnc.TrainFromJSONFile(path, cfg)).To(HaveOccurred())
			})
		})
	})

	Describe("TrainFromCSVFile", func() {
		BeforeEach(func() {
			path = "../datasets/shorebirds.csv"
		})

		It("Does not return an error", func() {
			Expect(knnc.TrainFromCSVFile(path, cfg)).NotTo(HaveOccurred())
		})

		When("It cannot read the data file", func() {
			BeforeEach(func() {
				path = "../datasets/thisdoesnotexist.csv"
			})

			It("Returns an error", func() {
				Expect(knnc.TrainFromCSVFile(path, cfg)).To(HaveOccurred())
			})
		})
	})

	Describe("Retrain", func() {
		var (
			trainingDS1, testDS1 *classifiers.DataSet
		)

		BeforeEach(func() {
			path = "../datasets/shorebirds.csv"
		})

		JustBeforeEach(func() {
			knnc.TrainFromCSVFile(path, cfg)
			trainingDS1 = knnc.TrainingData
			testDS1 = knnc.TestingData
		})

		When("The config is nil", func() {
			BeforeEach(func() {
				Expect(cfg).To(BeNil())
			})

			It("Replaces the existing training data with a different subset of the records", func() {
				knnc.Retrain(cfg)
				Expect(knnc.TrainingData.Records).To(HaveLen(len(trainingDS1.Records)))
				Expect(knnc.TrainingData.Records).NotTo(ConsistOf(trainingDS1.Records))
				Expect(knnc.TrainingData.Records).NotTo(ContainElements(trainingDS1.Records))
				Expect(trainingDS1.Records).NotTo(ContainElements(knnc.TrainingData.Records))
			})

			It("Replaces the existing testing data with a different subset of the records", func() {
				knnc.Retrain(cfg)
				Expect(knnc.TestingData.Records).To(HaveLen(len(testDS1.Records)))
				Expect(knnc.TestingData.Records).NotTo(ConsistOf(testDS1.Records))
				Expect(knnc.TestingData.Records).NotTo(ContainElements(testDS1.Records))
				Expect(testDS1.Records).NotTo(ContainElements(knnc.TestingData.Records))
			})
		})
	})

	Describe("Test", func() {
		BeforeEach(func() {
			path = "../fixtures/students.csv"
			cfg = &classifiers.DataSplitConfig{
				Method: classifiers.SplitSequential,
			}
		})

		When("The model has been trained", func() {
			JustBeforeEach(func() {
				err := knnc.TrainFromCSVFile(path, cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Does not return an error", func() {
				_, err := knnc.Test()
				Expect(err).NotTo(HaveOccurred())
			})

			It("Returns a result set for the test data", func() {
				results, _ := knnc.Test()
				resultRecords := make([]classifiers.Record, 0)
				for _, res := range results {
					resultRecords = append(resultRecords, res.Record)
				}
				Expect(resultRecords).To(ConsistOf(knnc.TestingData.Records))
			})

			It("Makes predictions about the classes of the testing records", func() {
				results, _ := knnc.Test()
				for _, res := range results {
					Expect(res.Predicted).NotTo(Equal(classifiers.NO_PREDICTION))
				}
			})

			It("Predicts the class based on the k nearest neighbors", func() {
				results, _ := knnc.Test()
				a := results.Analyze()
				// This is valid for the test data from students.csv
				Expect(a.Accuracy).To(Equal(1.0))
			})
		})

		When("The model has not been trained", func() {
			It("Returns an error", func() {
				_, err := knnc.Test()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
