package generator

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ScarletTanager/basilisk/classifiers"
)

var (
	configPath                       string
	outputPath                       string
)

type DataSetAttribute struct {
	Name                  string    `json:"name"`
	LowerBound            float64   `json:"lowerBound"`
	UpperBound            float64   `json:"upperBound"`
	AllocationsByQuintile []float64 `json:"allocationsByQuintile"`
}

type DatasetConfig struct {
	Classes    map[string][]DataSetAttribute `json:"classes"`
	EntryCount uint64                        `json:"entryCount"`
}

func init() {
	flag.StringVar(&configPath, "config", "config.json", "Path to configuration file in JSON")
	flag.StringVar(&outputPath, "output", "output.json", "Path to output file")
}

func main() {
	flag.Parse()

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
	for k, v := range datasetConfig.Classes {
		classNames = append(classNames, k)
	}

	attributeNames := make([]string, len(datasetConfig.Classes[classNames[0]]))
	for i, attr := range datasetConfig.Classes[classNames[0]] {
		attributeNames[i] = attr.Name
	}

	// Generate the data
	// This is a crappy way to do this if the number of records is really large
	records := make([]classifiers.Record, datasetConfig.EntryCount)
	for i := 0
}
