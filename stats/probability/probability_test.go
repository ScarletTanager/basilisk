package probability_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/stats/probability"
)

var _ = Describe("Probability", func() {

	Describe("MassDiscrete", func() {
		var (
			values []int
		)

		BeforeEach(func() {
			values = []int{3, 3, 1, 2, 3, 1, 1, 2, 3, 1}
		})

		It("Returns a correct pmf over the sample space", func() {
			pmf := probability.MassDiscrete(values)

			Expect(pmf(1)).To(Equal(0.4))
			Expect(pmf(2)).To(Equal(0.2))
			Expect(pmf(3)).To(Equal(0.4))

			totalProbability := float64(0)
			for _, v := range []int{1, 2, 3} {
				totalProbability += pmf(v)
			}

			Expect(totalProbability).To(Equal(1.0))
		})
	})
})
