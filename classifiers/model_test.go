package classifiers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/classifiers"
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
})
