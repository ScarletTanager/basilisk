package classifiers_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/wyvern"
)

var _ = Describe("Classifiers", func() {
	Describe("DataSet", func() {
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
		})
	})
})