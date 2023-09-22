package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/wyvern"
)

var (
	configPath   string
	outputPath   string
	outputFormat string
)

const (
	format_JSON = "json"
	format_CSV  = "csv"
)

type DataSetAttribute struct {
	Name                  string    `json:"name"`
	LowerBound            float64   `json:"lower"`
	UpperBound            float64   `json:"upper"`
	AllocationsByQuintile []float64 `json:"allocationsByQuintile"`
}

type DatasetConfig struct {
	Classes     map[string][]DataSetAttribute `json:"classes"`
	RecordCount int                           `json:"recordCount"`
}

func init() {
	flag.StringVar(&configPath, "config", "config.json", "Path to configuration file in JSON")
	flag.StringVar(&outputPath, "output", "output.json", "Path to output file")
	flag.StringVar(&outputFormat, "format", format_JSON, "Output format (default is JSON); csv and json are supported, value is case-insensitive")
}

func main() {
	flag.Parse()

	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case format_CSV:
	case format_JSON:
	default:
		log.Fatal(fmt.Sprintf("%s is not a valid format.  Supported values are 'csv' and 'json', case-insensitive.", outputFormat))
	}

	if configPath == "" {
		fmt.Fprintln(os.Stderr, "Must pass in path to config file with --config <path>")
		os.Exit(1)
	}

	jsonBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Errorf("While trying to read config file: %w", err))
	}

	var datasetConfig DatasetConfig

	err = json.Unmarshal(jsonBytes, &datasetConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing config file: %s\n", err)
		os.Exit(1)
	}

	classNames := make([]string, 0)
	for k, _ := range datasetConfig.Classes {
		classNames = append(classNames, k)
	}

	attributeNames := make([]string, len(datasetConfig.Classes[classNames[0]]))
	for i, attr := range datasetConfig.Classes[classNames[0]] {
		attributeNames[i] = attr.Name
	}

	// Generate the data
	// This is a crappy way to do this if the number of records is really large
	records := make([]classifiers.Record, datasetConfig.RecordCount)
	for ri := range records {
		records[ri] = classifiers.Record{
			AttributeValues: make(wyvern.Vector[float64], len(attributeNames)),
		}
	}

	// Assign each record index a random class
	recordIndicesByClass := assignClasses(datasetConfig.RecordCount, len(classNames))

	// Compute the attribute values for each class
	// recordsInClass is []int - values are indices into the _overall_ list of records from 0 to datasetConfig.RecordCount
	for classIdx, recordsInClass := range recordIndicesByClass {
		// Populate the class in the overall record table
		// TODO: we could do this more cleanly by having a single method
		// generate the list of records with the classes assigned
		for _, ri := range recordsInClass {
			records[ri].Class = classIdx
		}

		// Iterate over attributes on a per-class basis
		for attrIdx, attr := range datasetConfig.Classes[classNames[classIdx]] {
			// For each attribute, compute a quintile distribution for the current class
			// Assign records to quintiles according to the computed distribution
			// Each bucket is []int, values are the original indices into the overall record set (values of recordsInClass)
			countsByQuintile := computeQuintileDistribution(len(recordsInClass), attr.AllocationsByQuintile)
			recordsByQuintile := assignQuintiles(recordsInClass, countsByQuintile)
			for qi, quintile := range recordsByQuintile {
				for _, r := range quintile {
					val := computeAttributeValue(attr.LowerBound, attr.UpperBound, qi)
					records[r].AttributeValues[attrIdx] = val
				}
			}
		}
	}

	dataset, err := classifiers.NewDataSet(classNames, attributeNames, records)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating dataset: %s\n", err)
		os.Exit(1)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file: %s\n", err)
		os.Exit(1)
	}

	var datasetBytes []byte

	switch outputFormat {
	case format_JSON:
		datasetBytes, err = json.Marshal(dataset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling dataset: %s\n", err)
			os.Exit(1)
		}
	case format_CSV:
		datasetBytes = dataset.MarshalCSV()
	}

	_, err = f.Write(datasetBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output data: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

// Returns slice of slice(int) - outer indices are class indices, inner indices
// range up to the count of records assigned each class, inner values
// are the actual record indices assigned to each class
func assignClasses(recordCount, classCount int) [][]int {
	recordIndicesByClass := make([][]int, classCount)
	for i, _ := range recordIndicesByClass {
		recordIndicesByClass[i] = make([]int, 0)
	}

	for i := 0; i < recordCount; i++ {
		assignedClass := rand.Intn(classCount)
		recordIndicesByClass[assignedClass] = append(recordIndicesByClass[assignedClass], i)
	}

	return recordIndicesByClass
}

// Given a set of record indices ([]int), assign them to quintiles according to quintileCounts
func assignQuintiles(records []int, quintileCounts []int) [][]int {
	quintileRecords := make([][]int, 5)
	for i := range quintileRecords {
		quintileRecords[i] = make([]int, 0)
	}

	bucket := 0
	for ri := 0; ri < len(records); ri++ {
		// Yeah, risky
		// Find the first quintile which is not at capacity, starting with the one
		// to which we most recently added.
		for b := bucket; ; b++ {
			quintileIdx := b % 5
			if len(quintileRecords[quintileIdx]) < quintileCounts[quintileIdx] {
				quintileRecords[quintileIdx] = append(quintileRecords[quintileIdx], records[ri])
				bucket = b
				break
			}
		}
	}

	return quintileRecords
}

func computeQuintileDistribution(recordCount int, allocations []float64) []int {
	recordCountsByQuintile := make([]int, 5)
	for i, a := range allocations {
		recordCountsByQuintile[i] = int((a / 100) * float64(recordCount))
	}

	tot := 0
	for _, rc := range recordCountsByQuintile {
		tot += rc
	}

	if tot < recordCount {
		for i := 0; i < (recordCount - tot); i++ {
			recordCountsByQuintile[i%5] += 1
		}
	} else if tot > recordCount {
		for i := 0; i < (tot - recordCount); i++ {
			recordCountsByQuintile[i%5] -= 1
		}
	}

	return recordCountsByQuintile
}

func computeAttributeValue(lower, upper float64, quintile int) float64 {
	return ((upper - lower) * (float64((quintile*20)+rand.Intn(21)) / 100.0)) + lower
}
