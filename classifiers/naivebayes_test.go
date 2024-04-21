package classifiers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/sphinx/probability"
	"github.com/ScarletTanager/wyvern"
)

var _ = FDescribe("Naivebayes", func() {
	var (
		nbc  *classifiers.NaiveBayesClassifier
		cfg  *classifiers.DataSplitConfig
		path string
	)

	BeforeEach(func() {
		nbc = classifiers.NewNaiveBayes()
		cfg = nil
	})

	Context("TrainFromCSVFile", func() {
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

	Describe("ComputeClassPriors", func() {
		var (
			classPriors []float64
			classNames  []string
			records     []classifiers.Record
		)

		BeforeEach(func() {
			classNames = []string{"foo", "bar"}
			records = []classifiers.Record{
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{5.0, 5.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{13.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{20.0, 7.5},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{9.8, 10.0},
				},
			}
		})

		It("Computes the correct class priors", func() {
			classPriors = classifiers.ComputeClassPriors(classNames, records)
			Expect(classPriors).To(HaveLen(len(classNames)))
			Expect(classPriors[0]).To(Equal(0.6))
			Expect(classPriors[1]).To(Equal(0.4))
		})
	})

	Describe("DiscretizeAttributes", func() {
		var (
			attributeNames []string
			records        []classifiers.Record
		)

		BeforeEach(func() {
			attributeNames = []string{"a", "b"}
			records = []classifiers.Record{
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{5.0, 5.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{13.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{20.0, 7.5},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{9.8, 10.0},
				},
			}
		})

		It("Creates a set of intervals for each attribute", func() {
			attributeIntervals := classifiers.DiscretizeAttributes(attributeNames, records)
			Expect(attributeIntervals).To(HaveLen(2))
		})

		It("Creates intervals using the default config (10, equal size)", func() {
			attributeIntervals := classifiers.DiscretizeAttributes(attributeNames, records)
			Expect(attributeIntervals[0]).To(HaveLen(10))

			Expect(attributeIntervals[0][0].Lower).To(Equal(1.0))
			Expect(attributeIntervals[0][0].IncludesUpper).To(BeFalse()) // Equal distribution would be true

			Expect(attributeIntervals[0][9].Upper).To(Equal(20.0))
			Expect(attributeIntervals[0][9].IncludesUpper).To(BeTrue()) // Equal distribution would be true

			Expect(attributeIntervals[1]).To(HaveLen(10))
			Expect(attributeIntervals[1][0].Lower).To(Equal(1.0))
			Expect(attributeIntervals[1][0].IncludesUpper).To(BeFalse()) // Equal distribution would be true

			Expect(attributeIntervals[1][9].Upper).To(Equal(10.0))
			Expect(attributeIntervals[1][9].IncludesUpper).To(BeTrue()) // Equal distribution would be true
		})
	})

	Describe("GenerateClassAttributePosteriors", func() {
		var (
			classNames, attributeNames []string
			classAttrIntervals         []probability.Intervals
			records                    []classifiers.Record
		)

		BeforeEach(func() {
			classNames = []string{"a", "b"}
			attributeNames = []string{"width", "length"}
			records = []classifiers.Record{
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{1.0, 0.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{1.2, 0.1},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{2.3, 0.2},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{2.35, 0.3},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{2.4, 1.0},
				},
				// 5 records above
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{3.7, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{4.3, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{4.7, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{5.7, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{6.0, 2.0},
				},
				// 10 records above
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				// 15 records above
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 10.0},
				},
				// 20 records above
			}
		})

		JustBeforeEach(func() {
			classAttrIntervals = classifiers.DiscretizeAttributes(attributeNames, records)
		})

		It("Creates posteriors for each class, attribute, interval", func() {
			caps := classifiers.GenerateClassAttributePosteriors(classNames, attributeNames, classAttrIntervals, records)
			Expect(caps).To(HaveLen(2))
			Expect(caps[0]).To(HaveLen(len(attributeNames)))
			Expect(caps[1]).To(HaveLen(len(attributeNames)))
			for i := 0; i < len(caps); i++ {
				for j := 0; j < len(attributeNames); j++ {
					// We expect to have a probability value for each interval
					Expect(caps[i][j]).To(HaveLen(10))
				}
			}
		})

		It("Creates posteriors conditioned by the class", func() {
			caps := classifiers.GenerateClassAttributePosteriors(classNames, attributeNames, classAttrIntervals, records)
			Expect(caps[0][0][0]).NotTo(Equal(caps[1][0][0]))
		})

		It("Computes the correct posteriors", func() {
			caps := classifiers.GenerateClassAttributePosteriors(classNames, attributeNames, classAttrIntervals, records)

			Expect(caps[0][0][0]).To(Equal(0.2))
			Expect(caps[0][0][1]).To(Equal(0.0))
			Expect(caps[0][0][2]).To(Equal(0.3))
			Expect(caps[0][0][3]).To(Equal(0.0))
			Expect(caps[0][0][4]).To(Equal(0.0))
			Expect(caps[0][0][5]).To(Equal(0.1))
			Expect(caps[0][0][6]).To(Equal(0.1))
			Expect(caps[0][0][7]).To(Equal(0.1))
			Expect(caps[0][0][8]).To(Equal(0.0))
			Expect(caps[0][0][9]).To(Equal(0.2))

			Expect(caps[0][1][0]).To(Equal(0.4))
			Expect(caps[0][1][1]).To(Equal(0.1))
			Expect(caps[0][1][2]).To(Equal(0.5))
			Expect(caps[0][1][3]).To(Equal(0.0))
			Expect(caps[0][1][4]).To(Equal(0.0))
			Expect(caps[0][1][5]).To(Equal(0.0))
			Expect(caps[0][1][6]).To(Equal(0.0))
			Expect(caps[0][1][7]).To(Equal(0.0))
			Expect(caps[0][1][8]).To(Equal(0.0))
			Expect(caps[0][1][9]).To(Equal(0.0))

			Expect(caps[1][0][0]).To(Equal(1.0))
			Expect(caps[1][0][1]).To(Equal(0.0))
			Expect(caps[1][0][2]).To(Equal(0.0))
			Expect(caps[1][0][3]).To(Equal(0.0))
			Expect(caps[1][0][4]).To(Equal(0.0))
			Expect(caps[1][0][5]).To(Equal(0.0))
			Expect(caps[1][0][6]).To(Equal(0.0))
			Expect(caps[1][0][7]).To(Equal(0.0))
			Expect(caps[1][0][8]).To(Equal(0.0))
			Expect(caps[1][0][9]).To(Equal(0.0))

			Expect(caps[1][1][0]).To(Equal(0.0))
			Expect(caps[1][1][1]).To(Equal(0.9))
			Expect(caps[1][1][2]).To(Equal(0.0))
			Expect(caps[1][1][3]).To(Equal(0.0))
			Expect(caps[1][1][4]).To(Equal(0.0))
			Expect(caps[1][1][5]).To(Equal(0.0))
			Expect(caps[1][1][6]).To(Equal(0.0))
			Expect(caps[1][1][7]).To(Equal(0.0))
			Expect(caps[1][1][8]).To(Equal(0.0))
			Expect(caps[1][1][9]).To(Equal(0.1))
		})
	})

	Describe("GenerateClassConditionedPosteriors", func() {
		var (
			classNames, attributeNames []string
			records                    []classifiers.Record
			intervals                  []probability.Intervals
			caps                       classifiers.ClassAttributePosteriors
		)

		BeforeEach(func() {
			classNames = []string{"a", "b"}
			attributeNames = []string{"width", "length"}
			records = []classifiers.Record{
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{1.0, 0.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{1.2, 0.1},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{2.3, 0.2},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{2.35, 0.3},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{2.4, 1.0},
				},
				// 5 records above
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{3.7, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{4.3, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{4.7, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{5.7, 2.0},
				},
				{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{6.0, 2.0},
				},
				// 10 records above
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				// 15 records above
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 1.0},
				},
				{
					Class:           1,
					AttributeValues: wyvern.Vector[float64]{1.0, 10.0},
				},
				// 20 records above
			}
		})

		JustBeforeEach(func() {
			intervals = classifiers.DiscretizeAttributes(attributeNames, records)
			caps = classifiers.GenerateClassAttributePosteriors(classNames, attributeNames, intervals, records)
		})

		// Skimping a bit here on the tests
		It("Computes correct values for each class and vector", func() {
			ccps := classifiers.GenerateClassConditionedPosteriors(classNames,
				attributeNames,
				intervals,
				caps)

			Expect(ccps).To(HaveLen(2))
			Expect(ccps[0]).To(HaveLen(100))
			Expect(ccps[1]).To(HaveLen(100))

			// With attributes indexed 0 and 1, each with 10 intervals, vector index 15
			// has probability P(A1==1)P(A2==5).

			for classIdx, _ := range classNames {
				for tens := 0; tens < 10; tens++ {
					for ones := 0; ones < 10; ones++ {
						expected := caps[classIdx][0][tens] * caps[classIdx][1][ones]
						Expect(ccps[classIdx][(tens*10)+ones]).To(Equal(expected))
					}
				}
			}
		})
	})

	Describe("NaiveBayesClassifier", func() {
		var (
			nbc *classifiers.NaiveBayesClassifier
		)

		BeforeEach(func() {
			nbc = classifiers.NewNaiveBayes()
		})

		JustBeforeEach(func() {
			Expect(nbc.VectorConditionedClassProbabilities).To(BeNil())
		})

		Describe("TrainFromDataset", func() {
			var (
				dataset *classifiers.DataSet
				err     error
			)

			BeforeEach(func() {
				dataset, err = classifiers.FromJSONFile("../datasets/widgets.json")
				Expect(err).NotTo(HaveOccurred())
			})

			When("The DataSet is valid", func() {
				It("Trains the model with vector-conditioned class posteriors", func() {
					nbc.TrainFromDataset(dataset, &classifiers.DataSplitConfig{
						TrainingShare: 0.80,
						Method:        classifiers.SplitSequential,
					})
					Expect(nbc.VectorConditionedClassProbabilities).NotTo(BeNil())
				})
			})
		})
	})
})
