package classifiers

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
	Predicted int
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
