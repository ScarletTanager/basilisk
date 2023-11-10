package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ScarletTanager/basilisk/dsgen"
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

	var datasetConfig dsgen.DatasetConfig

	err = json.Unmarshal(jsonBytes, &datasetConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing config file: %s\n", err)
		os.Exit(1)
	}

	dataset, err := dsgen.GenerateDataset(&datasetConfig)
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
