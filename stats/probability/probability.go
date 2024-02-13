package probability

// MassDiscrete returns a probability mass function (PMF) over the range of values
// which can be assigned to a given random variable.  The variable has a discrete
// range, and the values must represent the full sample space.
func MassDiscrete(values []int) func(int) float64 {
	count := float64(len(values))
	if count == 0 {
		return func(int) float64 {
			return 0
		}
	}

	valCounts := make(map[int]float64)
	for _, val := range values {
		valCounts[val] += 1.0
	}

	return func(x int) float64 {
		return valCounts[x] / count
	}
}
