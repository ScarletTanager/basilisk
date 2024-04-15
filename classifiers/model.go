package classifiers

import (
	"math"

	"github.com/ScarletTanager/wyvern"
)

const (
	NO_PREDICTION = -1
)

// Classifier is the interface implemented by all classification models
type Classifier interface {
	TrainFromCSV([]byte, *DataSplitConfig) error
	TrainFromCSVFile(string, *DataSplitConfig) error
	TrainFromDataset(*DataSet, *DataSplitConfig) error
	TrainFromJSON([]byte, *DataSplitConfig) error
	TrainFromJSONFile(string, *DataSplitConfig) error
	Retrain(*DataSplitConfig) error
	Test() (TestResults, error)
	Type() string
	Data() (*DataSet, *DataSet)
	Config() interface{}
}

type TestResults []TestResult

type TestResult struct {
	Record
	Predicted   int
	Probability float64
	// Votes is the number of votes (in a nearest neighbors model) for the predicated class
	Votes int
}

type TestResultsAnalysis struct {
	ResultCount    int     `json:"results"`
	CorrectCount   int     `json:"correct"`
	IncorrectCount int     `json:"incorrect"`
	Accuracy       float64 `json:"accuracy"`
}

func (trs TestResults) Analyze() TestResultsAnalysis {
	analysis := TestResultsAnalysis{
		ResultCount: len(trs),
	}

	for _, result := range trs {
		if result.Class == result.Predicted {
			analysis.CorrectCount++
		} else {
			analysis.IncorrectCount++
		}
	}

	analysis.Accuracy = float64(analysis.CorrectCount) / float64(analysis.ResultCount)

	return analysis
}

const (
	DistanceMethod_Euclidean = "euclidean"
	DistanceMethod_Manhattan = "manhattan"
)

type DistanceFunction func(wyvern.Vector[float64], wyvern.Vector[float64]) float64

// For both distance functions, we're assuming that the vectors have the same
// dimensionality (number of components).  Don't use these with vectors that
// don't have the same dimensionality and expect things to "just work."

// EuclideanDistance is (for now) just a convenience method to return
// euclidean distance - we will probably change the implementation if/when
// we support attribute types other than float64
func EuclideanDistance(a, b wyvern.Vector[float64]) float64 {
	return a.Difference(b).Magnitude()
}

// ManhattanDistance returns the Manhattan (city block) distance
// between two points (represented as vectors).
func ManhattanDistance(a, b wyvern.Vector[float64]) float64 {
	var distance float64
	for _, component := range a.Difference(b) {
		// Remember that we want the magnitude of the difference
		distance += math.Abs(component)
	}

	return distance
}
