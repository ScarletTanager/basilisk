package classifiers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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

// FromCSV builds a DataSet from CSV data.  Returns nil and an error if the data cannot be processed correctly.
func FromCSV(dsCsv []byte) (*DataSet, error) {
	// This is not memory-efficient, since it slurps in the entire slice of data.
	// TODO: process this from a bytes.Reader or similar?
	lines := bytes.Split(dsCsv, []byte("\n"))

	headerFields := strings.Split(string(lines[0]), ",")
	attributeNames := headerFields[:len(headerFields)-1]
	attributeCount := len(attributeNames)

	classNameMap := make(map[string]int)
	classIdx := 0

	records := make([]Record, 0)

	for lineNo, line := range lines[1:] {
		lineFields := bytes.Split(line, []byte(","))

		// Check that we have the correct number of attributes
		attributeValsRaw := lineFields[:len(lineFields)-1]
		if len(attributeValsRaw) != attributeCount {
			// Are we at the last line, and this is a blank line?
			if lineNo == len(lines)-1 || len(line) == 0 {
				break
			}

			// Add two to start at 1, header line is line 1 (so we skipped it)
			return nil, fmt.Errorf("Invalid data at line %d", lineNo+2)
		}

		// Parse the attribute columns
		attributeValues := make(wyvern.Vector[float64], attributeCount)
		for attrIdx, attrValRaw := range attributeValsRaw {
			if attrValue, conversionErr := strconv.ParseFloat(string(attrValRaw), 64); conversionErr != nil {
				return nil, fmt.Errorf("Unable to parse attribute value %s, index %d, at line %d into float64", attrValRaw, attrIdx, lineNo+2)
			} else {
				attributeValues[attrIdx] = attrValue
			}
		}

		rec := Record{
			AttributeValues: attributeValues,
		}

		// If we have not seen the className before, add it to the map and bump the index
		// for the next value
		className := string(lineFields[len(lineFields)-1])
		if _, ok := classNameMap[className]; !ok {
			classNameMap[className] = classIdx
			rec.Class = classIdx
			classIdx++
		}

		records = append(records, rec)
	}

	// Convert the class names to a slice
	classNames := make([]string, len(classNameMap))
	for name, idx := range classNameMap {
		classNames[idx] = name
	}

	return NewDataSet(classNames, attributeNames, records)
}

// FromCSVFile reads the CSV-formatted file and creates a DataSet from it.
func FromCSVFile(path string) (*DataSet, error) {
	csvBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("When attempting to read CSV file: %w", err)
	}

	return FromCSV(csvBytes)
}

// MarshalCSV converts the DataSet to a byte slice containing the CSV representation (including a header
// row listing the attributes and terminated by the column header for the class column).
func (ds *DataSet) MarshalCSV() []byte {
	var buf bytes.Buffer

	// Write out the attribute names in the header
	for _, attrName := range ds.AttributeNames {
		buf.WriteString(attrName + ",")
	}

	// Finish the header
	buf.WriteString("class\n")

	for _, rec := range ds.Records {
		for _, val := range rec.AttributeValues {
			buf.WriteString(fmt.Sprintf("%f,", val))
		}
		buf.WriteString(fmt.Sprintf("%s\n", ds.ClassNames[rec.Class]))
	}

	return buf.Bytes()
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

// Split divides the dataset into separate datasets - the first is the training
// data, the second is the test data
// Passing nil for the config results in a random split with 75% of the records used for training.
func (ds *DataSet) Split(cfg *DataSplitConfig) (*DataSet, *DataSet, error) {
	return nil, nil, nil
}

type DataSplitMethod int

const (
	SplitRandom DataSplitMethod = iota
	SplitSequential
)

type DataSplitConfig struct {
	TrainingShare float64
	Method        DataSplitMethod
}

type TestResults []TestResult

type TestResult struct {
	Record
	Predicted int
}
