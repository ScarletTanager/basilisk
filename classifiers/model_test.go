package classifiers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/wyvern"
)

var _ = Describe("Model", func() {
	Describe("TestResults", func() {
		var (
			results classifiers.TestResults
		)

		Describe("Analyze", func() {
			BeforeEach(func() {
				ds, err := classifiers.FromCSVFile("../fixtures/students.csv")
				Expect(err).NotTo(HaveOccurred())
				results = make(classifiers.TestResults, len(ds.Records))
				for i, r := range ds.Records {
					results[i] = classifiers.TestResult{
						Record:    r,
						Predicted: r.Class,
					}
				}
			})

			It("Produces an accurate analysis", func() {
				analysis := results.Analyze()
				Expect(analysis.ResultCount).To(Equal(len(results)))
				Expect(analysis.CorrectCount).To(Equal(len(results)))
				Expect(analysis.IncorrectCount).To(Equal(0))
				Expect(analysis.Accuracy).To(Equal(1.0))
			})

			When("Some predications were incorrect", func() {
				BeforeEach(func() {
					for i, r := range results[9:] {
						results[i].Predicted = r.Class + 1
					}
				})

				It("Produces an accurate analysis", func() {
					analysis := results.Analyze()
					Expect(analysis.ResultCount).To(Equal(len(results)))
					Expect(analysis.CorrectCount).To(Equal(len(results) - 3))
					Expect(analysis.IncorrectCount).To(Equal(3))
					Expect(analysis.Accuracy).To(Equal(.75))
				})
			})
		})
	})

	Describe("Distance computation", func() {
		var (
			a, b             wyvern.Vector[float64]
			expectedDistance float64
		)

		BeforeEach(func() {
			a = wyvern.Vector[float64]{
				5.5, 12.0, 6, -18.3,
			}

			b = wyvern.Vector[float64]{
				-3.2, 18.0, 5.5, -27.0,
			}
		})

		JustBeforeEach(func() {
			Expect(a).To(HaveLen(len(b)))
		})

		Describe("EuclideanDistance", func() {
			JustBeforeEach(func() {
				expectedDistance = a.Difference(b).Magnitude()
			})

			It("Returns the magnitude of the vector difference", func() {
				Expect(classifiers.EuclideanDistance(a, b)).To(Equal(expectedDistance))
			})
		})

		Describe("ManhattanDistance", func() {
			JustBeforeEach(func() {
				expectedDistance = 0
				for i, v := range a {
					expectedDistance += (v - b[i])
				}
			})

			It("Returns the city block distance between the two points", func() {
				Expect(classifiers.ManhattanDistance(a, b)).To(Equal(expectedDistance))
			})
		})
	})
})
