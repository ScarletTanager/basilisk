package classifiers

const (
	NO_PREDICTION = -1
)

type TestResults []TestResult

type TestResult struct {
	Record
	Predicted int
}

type TestResultsAnalysis struct {
	ResultCount    int
	CorrectCount   int
	IncorrectCount int
	Accuracy       float64
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
