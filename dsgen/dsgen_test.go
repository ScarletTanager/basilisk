package dsgen_test

import (
	"github.com/ScarletTanager/basilisk/dsgen"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dsgen", func() {
	var (
		cfg     *dsgen.DatasetConfig
		count   int
		classes map[string][]dsgen.DataSetAttribute
	)

	JustBeforeEach(func() {
		cfg = &dsgen.DatasetConfig{
			RecordCount: count,
			Classes:     classes,
		}
	})

	Describe("DatasetConfig", func() {
		Describe("ClassNames", func() {
			When("The configuration defines classes", func() {
				BeforeEach(func() {
					classes = map[string][]dsgen.DataSetAttribute{
						"class1": nil,
						"class2": nil,
						"class3": nil,
					}
				})

				JustBeforeEach(func() {
					Expect(cfg.Classes).To(HaveLen(len(classes)))
				})

				It("Returns the names of all defined classes", func() {
					definedClassNames := make([]string, 0)
					for name, _ := range cfg.Classes {
						definedClassNames = append(definedClassNames, name)
					}
					Expect(cfg.ClassNames()).To(ConsistOf(definedClassNames))
				})
			})

			When("The configuration does not define any classes", func() {
				BeforeEach(func() {
					classes = nil
				})

				It("Returns an empty slice", func() {
					Expect(cfg.ClassNames()).To(BeEmpty())
				})
			})
		})
	})

	Describe("GenerateDataset", func() {
		BeforeEach(func() {
			count = 30
			classes = map[string][]dsgen.DataSetAttribute{
				"small": {
					{
						Name:                  "length",
						LowerBound:            0.0,
						UpperBound:            10.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
					{
						Name:                  "height",
						LowerBound:            0.0,
						UpperBound:            10.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
					{
						Name:                  "width",
						LowerBound:            0.0,
						UpperBound:            10.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
				},
				"medium": {
					{
						Name:                  "length",
						LowerBound:            10.0,
						UpperBound:            20.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
					{
						Name:                  "height",
						LowerBound:            10.0,
						UpperBound:            20.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
					{
						Name:                  "width",
						LowerBound:            10.0,
						UpperBound:            20.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
				},
				"large": {
					{
						Name:                  "length",
						LowerBound:            20.0,
						UpperBound:            50.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
					{
						Name:                  "height",
						LowerBound:            20.0,
						UpperBound:            50.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
					{
						Name:                  "width",
						LowerBound:            20.0,
						UpperBound:            50.0,
						AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
					},
				},
			}
		})

		When("The config is valid", func() {
			It("Generates a dataset", func() {
				ds, _ := dsgen.GenerateDataset(cfg)
				Expect(ds).NotTo(BeNil())
			})

			It("Does not return an error", func() {
				_, err := dsgen.GenerateDataset(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Generates a dataset with the correct number of records", func() {
				ds, _ := dsgen.GenerateDataset(cfg)
				Expect(ds.Records).NotTo(BeNil())
				Expect(ds.Records).To(HaveLen(count))
			})

			It("Generates a dataset including all of the specified classes", func() {
				ds, _ := dsgen.GenerateDataset(cfg)
				Expect(ds.ClassNames).To(ConsistOf(cfg.ClassNames()))
			})

			It("Generates a dataset including all of the defined attributes", func() {
				ds, _ := dsgen.GenerateDataset(cfg)
				Expect(ds.AttributeNames).To(ConsistOf("length", "height", "width"))
			})
		})
	})

	Describe("AssignQuintiles", func() {
		var (
			indices        []int
			quintileCounts []int
			allocations    []float64
		)

		BeforeEach(func() {
			count = 1000
			allocations = []float64{35.0, 13.0, 12.0, 20.0, 20.0}
		})

		JustBeforeEach(func() {
			indices = make([]int, count)
			for i := 0; i < count; i++ {
				indices[i] = i
			}

			quintileCounts, _ = dsgen.ComputeQuintileDistribution(count, allocations)
		})

		It("Assigns each quintile a number of indices according to the distributions", func() {
			indicesByQuintile, _ := dsgen.AssignQuintiles(indices, quintileCounts)
			Expect(indicesByQuintile).To(HaveLen(5))
			for i, q := range indicesByQuintile {
				Expect(q).To(HaveLen(quintileCounts[i]))
			}
		})

		When("More than 5 'quintiles' are specified", func() {
			JustBeforeEach(func() {
				quintileCounts = append(quintileCounts, 10)
			})

			It("Returns an error", func() {
				_, err := dsgen.AssignQuintiles(indices, quintileCounts)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ComputeAttributeValue", func() {
		var (
			lower, upper float64
		)

		BeforeEach(func() {
			count = 1000
			lower = 0.0
			upper = 100.0
		})

		It("Generates a random attribute value within the correct quintile", func() {
			for i := 0; i < count; i++ {
				for qi := 0; qi < 5; qi++ {
					if qi > 0 {
						Expect(dsgen.ComputeAttributeValue(lower, upper, qi)).To(SatisfyAll(
							BeNumerically(">", float64(qi)*(upper/5.0)),
							BeNumerically("<=", float64(qi+1)*(upper/5.0))))
					} else {
						Expect(dsgen.ComputeAttributeValue(lower, upper, qi)).To(SatisfyAll(
							BeNumerically(">=", float64(qi)*(upper/5.0)),
							BeNumerically("<=", float64(qi+1)*(upper/5.0))))
					}
				}
			}
		})
	})

	Describe("ComputeQuintileDistribution", func() {
		var (
			allocations []float64
		)

		BeforeEach(func() {
			allocations = make([]float64, 5)
		})

		When("The record count is positive", func() {
			BeforeEach(func() {
				count = 1000
			})

			When("The total allocations exceed 100", func() {
				BeforeEach(func() {
					allocations[0] = 30.0
					allocations[1] = 30.0
					allocations[2] = 20.0
					allocations[3] = 20.0
					allocations[4] = 0.1
				})

				It("Returns an error", func() {
					distributions, err := dsgen.ComputeQuintileDistribution(count, allocations)
					Expect(distributions).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})

			When("There are too many allocations (more than 5)", func() {
				BeforeEach(func() {
					allocations = append(allocations, 0.5)
				})

				It("Returns an error", func() {
					distributions, err := dsgen.ComputeQuintileDistribution(count, allocations)
					Expect(distributions).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})

			When("The allocations total to 100.0", func() {
				BeforeEach(func() {
					allocations = []float64{35.0, 13.0, 12.0, 20.0, 20.0}
				})

				It("Does not return an error", func() {
					_, err := dsgen.ComputeQuintileDistribution(count, allocations)
					Expect(err).NotTo(HaveOccurred())
				})

				It("Computes the quintile distributions according to the specified allocations", func() {
					distributions, _ := dsgen.ComputeQuintileDistribution(count, allocations)
					Expect(distributions).To(HaveLen(5))
					for i, d := range distributions {
						Expect(d).To(Equal(int(float64(count) * (allocations[i] / 100.0))))
					}
				})
			})
		})

		When("The record count is not positive", func() {
			BeforeEach(func() {
				count = 0
			})

			It("Returns an error", func() {
				_, err := dsgen.ComputeQuintileDistribution(count, allocations)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
