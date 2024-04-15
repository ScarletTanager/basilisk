package classifiers

type ClassifierImplementation struct {
	RawData      *DataSet
	TrainingData *DataSet
	TestingData  *DataSet
	Results      TestResults
}
