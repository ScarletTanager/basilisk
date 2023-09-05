package classifiers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ScarletTanager/wyvern"
)

type DataSet struct {
	ClassNames     []string `json:"classes"`
	AttributeNames []string `json:"attributes"`
	Records        []Record `json:"data"`
}

type Record struct {
	Class           int                    `json:"class"`
	AttributeValues wyvern.Vector[float64] `json:"values"`
}

func NewDataSet(classes, attributes []string, data []Record) (*DataSet, error) {
	if data != nil {
		if !dataHasValidClasses(classes, data) {
			return nil, errors.New("At least one record has an invalid class")
		}

		if !dataHasValidAttributes(attributes, data) {
			return nil, errors.New("At least one record has too many attributes")
		}
	}

	return &DataSet{
		ClassNames:     classes,
		AttributeNames: attributes,
		Records:        data,
	}, nil
}

func FromJSON(dsJson []byte) (*DataSet, error) {
	var ds DataSet
	err := json.Unmarshal(dsJson, &ds)
	if err != nil {
		return nil, fmt.Errorf("While creating DataSet from JSON: %w", err)
	}

	return NewDataSet(ds.ClassNames, ds.AttributeNames, ds.Records)
}

func FromJSONFile(path string) (*DataSet, error) {
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("When attempting to read JSON from file: %w", err)
	}

	return FromJSON(jsonBytes)
}

func dataHasValidClasses(classes []string, data []Record) bool {
	for _, r := range data {
		if r.Class < 0 || r.Class >= len(classes) {
			return false
		}
	}
	return true
}

func dataHasValidAttributes(attributes []string, data []Record) bool {
	for _, r := range data {
		if len(r.AttributeValues) > len(attributes) {
			return false
		}
	}
	return true
}

func (ds *DataSet) Classes() []string {
	return ds.ClassNames
}

func (ds *DataSet) Attributes() []string {
	return ds.AttributeNames
}

type TestResults []TestResult

type TestResult struct {
	Record
	Predicted int
}
