package classifiers

type TestResults []TestResult

type TestResult struct {
	Record
	Predicted int
}
