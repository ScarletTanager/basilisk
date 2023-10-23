package classifiers_test

import (
	"encoding/json"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/wyvern"
)

var _ = Describe("DataSet", func() {
	var (
		ds             *classifiers.DataSet
		classes, attrs []string
		data           []classifiers.Record
	)

	BeforeEach(func() {
		classes = []string{
			"small",
			"medium",
			"large",
		}

		attrs = []string{
			"chest",
			"sleeve",
			"neck",
		}

		data = []classifiers.Record{
			{
				Class:           0,
				AttributeValues: wyvern.Vector[float64]{1, 1, 1},
			},
			{
				Class:           1,
				AttributeValues: wyvern.Vector[float64]{1, 1, 2},
			},
			{
				Class:           2,
				AttributeValues: wyvern.Vector[float64]{1, 2, 2},
			},
			{
				Class:           0,
				AttributeValues: wyvern.Vector[float64]{0, .5, 1},
			},
		}
	})

	Describe("NewDataSet", func() {
		It("Returns a pointer to the new DataSet and does not return an error", func() {
			dataSet, e := classifiers.NewDataSet(classes, attrs, data)
			Expect(e).NotTo(HaveOccurred())
			Expect(dataSet).NotTo(BeNil())
		})

		When("The data contains records with invalid classes", func() {
			BeforeEach(func() {
				data = append(data, classifiers.Record{
					Class:           3,
					AttributeValues: wyvern.Vector[float64]{4, 4, 4},
				})
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.NewDataSet(classes, attrs, data)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})

		When("The data contains records with too many attributes", func() {
			BeforeEach(func() {
				data = append(data, classifiers.Record{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{0, 1, 0, 1},
				})
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.NewDataSet(classes, attrs, data)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})
	})

	Describe("FromJSON", func() {
		var (
			dsJson []byte
		)

		JustBeforeEach(func() {
			ds = &classifiers.DataSet{
				ClassNames:     classes,
				AttributeNames: attrs,
				Records:        data,
			}
			var e error
			dsJson, e = json.Marshal(ds)
			Expect(e).NotTo(HaveOccurred())
		})

		It("Returns a non-nil DataSet pointer and does not return an error", func() {
			dataSet, e := classifiers.FromJSON(dsJson)
			Expect(dataSet).NotTo(BeNil())
			Expect(e).NotTo(HaveOccurred())
		})

		When("The JSON contains records with invalid classes", func() {
			BeforeEach(func() {
				data = append(data, classifiers.Record{
					Class:           3,
					AttributeValues: wyvern.Vector[float64]{4, 4, 4},
				})
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.FromJSON(dsJson)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})

		When("The JSON contains records with too many attributes", func() {
			BeforeEach(func() {
				data = append(data, classifiers.Record{
					Class:           0,
					AttributeValues: wyvern.Vector[float64]{0, 1, 0, 1},
				})
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.FromJSON(dsJson)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})
	})

	Describe("FromJSONFile", func() {
		var (
			path string
		)

		BeforeEach(func() {
			path = "../datasets/shorebirds.json"
		})

		It("Returns a pointer to a valid DataSet and does not return an error", func() {
			dataSet, e := classifiers.FromJSONFile(path)
			Expect(e).NotTo(HaveOccurred())
			Expect(dataSet).NotTo(BeNil())
		})

		When("The file cannot be opened", func() {
			BeforeEach(func() {
				path = "../datasets/thisdoesnotexist.json"
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.FromJSONFile(path)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})

		When("The file does not contain valid JSON", func() {
			BeforeEach(func() {
				path = "../fixtures/shorebirds_bad.json"
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.FromJSONFile(path)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})

		When("The JSON file contains records with invalid classes", func() {
			BeforeEach(func() {
				path = "../fixtures/shorebirds_badclasses.json"
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.FromJSONFile(path)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})

		When("The JSON file contains records with too many attributes", func() {
			BeforeEach(func() {
				path = "../fixtures/shorebirds_badattributes.json"
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.FromJSONFile(path)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})
	})

	Describe("FromCSV", func() {
		var (
			sourceCSV []byte
		)

		BeforeEach(func() {
			sourceDS, _ := classifiers.FromJSONFile("../datasets/shorebirds.json")
			sourceCSV = sourceDS.MarshalCSV()
		})

		It("Creates a valid DataSet from the source CSV", func() {
			ds, e := classifiers.FromCSV(sourceCSV)
			Expect(e).NotTo(HaveOccurred())
			Expect(ds).NotTo(BeNil())
		})

		When("The CSV contains records with non-float-valued attributes", func() {
			BeforeEach(func() {
				var err error
				sourceCSV, err = os.ReadFile("../fixtures/shorebirds_badvalues.csv")
				Expect(err).NotTo(HaveOccurred())
			})

			It("Returns nil and an error", func() {
				ds, e := classifiers.FromCSV(sourceCSV)
				Expect(e).To(HaveOccurred())
				Expect(ds).To(BeNil())
			})
		})

		When("The CSV contains records with the wrong number of attributes/columns", func() {
			BeforeEach(func() {
				var err error
				sourceCSV, err = os.ReadFile("../fixtures/shorebirds_badattributes.csv")
				Expect(err).NotTo(HaveOccurred())
			})

			It("Returns nil and an error", func() {
				ds, e := classifiers.FromCSV(sourceCSV)
				Expect(e).To(HaveOccurred())
				Expect(ds).To(BeNil())
			})
		})
	})

	Describe("FromCSVFile", func() {
		var (
			path string
		)

		BeforeEach(func() {
			path = "../datasets/shorebirds.csv"
		})

		It("Creates a valid DataSet from the file contents", func() {
			ds, e := classifiers.FromCSVFile(path)
			Expect(ds).NotTo(BeNil())
			Expect(e).NotTo(HaveOccurred())
		})

		When("The file cannot be opened", func() {
			BeforeEach(func() {
				path = "../datasets/thisdoesnotexist.json"
			})

			It("Returns nil and an error", func() {
				dataSet, e := classifiers.FromCSVFile(path)
				Expect(dataSet).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})

		When("The file contains records with the wrong number of columns/attributes", func() {
			BeforeEach(func() {
				path = "../fixtures/shorebirds_badattributes.csv"
			})

			It("Returns nil and an error", func() {
				ds, e := classifiers.FromCSVFile(path)
				Expect(ds).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})

		When("The file contains records with bad values (cannot be parsed into float64)", func() {
			BeforeEach(func() {
				path = "../fixtures/shorebirds_badvalues.csv"
			})

			It("Returns nil and an error", func() {
				ds, e := classifiers.FromCSVFile(path)
				Expect(ds).To(BeNil())
				Expect(e).To(HaveOccurred())
			})
		})
	})

	Describe("DataSet Methods", func() {
		JustBeforeEach(func() {
			ds, _ = classifiers.NewDataSet(classes, attrs, data)
			Expect(ds).NotTo(BeNil())
		})

		Describe("Classes", func() {
			It("Returns the class names", func() {
				Expect(ds.Classes()).To(Equal(classes))
			})
		})

		Describe("Attributes", func() {
			It("Returns the attribute names", func() {
				Expect(ds.Attributes()).To(Equal(attrs))
			})
		})

		Describe("MarshalCSV", func() {
			var (
				targetCSV []byte
			)

			BeforeEach(func() {
				targetCSV = []byte(`chest,sleeve,neck,class
1.000000,1.000000,1.000000,small
1.000000,1.000000,2.000000,medium
1.000000,2.000000,2.000000,large
0.000000,0.500000,1.000000,small
`)
			})

			It("Converts the dataset into a byte slice containing a valid CSV representation", func() {
				csv := ds.MarshalCSV()
				Expect(csv).To(Equal(targetCSV))
			})
		})

		Describe("Split", func() {
			var (
				splitCfg *classifiers.DataSplitConfig
				sourceDS *classifiers.DataSet
			)

			BeforeEach(func() {
				var e error
				splitCfg = nil
				sourceDS, e = classifiers.FromCSVFile("../datasets/shorebirds.csv")
				Expect(e).NotTo(HaveOccurred())
				Expect(sourceDS).NotTo(BeNil())
				// Expect(len(sourceDS.Records)).To(Equal(300))
			})

			It("Defaults to a 75/25 split, randomized", func() {
				testDataSetSplit(sourceDS, splitCfg)
			})

			When("A DataSplitConfig is specified", func() {
				BeforeEach(func() {
					splitCfg = &classifiers.DataSplitConfig{}
				})

				When("Without explicitly setting any config fields", func() {
					It("Uses the default 75/25 randomized split", func() {
						testDataSetSplit(sourceDS, splitCfg)
					})
				})

				When("With a training share specified", func() {
					BeforeEach(func() {
						splitCfg.TrainingShare = .9
					})
					It("Allocates the specified percentage to training data and the rest to test data, randomized", func() {
						testDataSetSplit(sourceDS, splitCfg)
					})
				})

				When("With a split method specified", func() {
					When("and the method is SplitRandom", func() {
						BeforeEach(func() {
							splitCfg.Method = classifiers.SplitRandom
						})

						It("Splits according to the specified method", func() {
							testDataSetSplit(sourceDS, splitCfg)
						})
					})

					When("and the method is SplitSequential", func() {
						BeforeEach(func() {
							splitCfg.Method = classifiers.SplitSequential
						})

						It("Splits according to the specified method", func() {
							testDataSetSplit(sourceDS, splitCfg)
						})
					})
				})

				When("With both training share and method specified", func() {
					BeforeEach(func() {
						splitCfg.TrainingShare = .65
						splitCfg.Method = classifiers.SplitSequential
					})

					It("Splits according to both the specified training share and method", func() {
						testDataSetSplit(sourceDS, splitCfg)
					})
				})
			})

			// It("Does not modify the original DataSet", func() {
			// })
		})
	})
})

func testDataSetSplit(ds *classifiers.DataSet, cfg *classifiers.DataSplitConfig) {
	var (
		trainingShare float64
		method        classifiers.DataSplitMethod
	)

	if cfg == nil {
		trainingShare = classifiers.DEFAULT_TRAINING_SHARE
		method = classifiers.SplitRandom
	} else {
		if cfg.TrainingShare == 0.0 {
			trainingShare = classifiers.DEFAULT_TRAINING_SHARE
		} else {
			trainingShare = cfg.TrainingShare
		}
		method = cfg.Method
	}

	trainingRecordCount := int(float64(len(ds.Records)) * trainingShare)
	testRecordCount := len(ds.Records) - trainingRecordCount

	trainingDS1, testDS1, err := ds.Split(cfg)
	Expect(trainingDS1).NotTo(BeNil())
	Expect(testDS1).NotTo(BeNil())
	Expect(err).NotTo(HaveOccurred())

	Expect(trainingDS1.Records).To(HaveLen(trainingRecordCount))
	Expect(testDS1.Records).To(HaveLen(testRecordCount))

	// Check that we have all the records and didn't throw any away
	allRecords := append(trainingDS1.Records, testDS1.Records...)
	Expect(allRecords).To(ConsistOf(ds.Records))

	if method == classifiers.SplitRandom {
		// Split a second time, we expect a different split as evidence of randomization
		trainingDS2, testDS2, _ := ds.Split(cfg)
		Expect(trainingDS2).NotTo(BeNil())
		Expect(testDS2).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(trainingDS2.Records).To(HaveLen(trainingRecordCount))
		Expect(testDS2.Records).To(HaveLen(testRecordCount))

		Expect(trainingDS2.Records).NotTo(ConsistOf(trainingDS1.Records))
		Expect(testDS2.Records).NotTo(ConsistOf(testDS1.Records))
	} else {
		// Sequential
		Expect(trainingDS1.Records).To(ConsistOf(ds.Records[:trainingRecordCount]))
		Expect(testDS1.Records).To(ConsistOf(ds.Records[trainingRecordCount:]))
	}
}
