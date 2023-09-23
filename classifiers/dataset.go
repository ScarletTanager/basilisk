package classifiers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/ScarletTanager/wyvern"
	"golang.org/x/exp/slices"
)

const (
	DEFAULT_TRAINING_SHARE = .75
)

var (
	maxRecordCount int
)

func init() {
	// maxRecordCount holds the cube root of math.MaxInt.  DataSets containg more
	// records than this have to be randomized in parts over multiple passes
	// (similar to shuffling two or more decks of cards into a single fully shuffled deck).
	maxRecordCount = int(math.Pow(float64(math.MaxInt), 1.0/3.0))
	// Make it an even number so we can use half of it
	if maxRecordCount%2 != 0 {
		maxRecordCount--
	}
}

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
// This does not modify the original DataSet.
func (ds *DataSet) Split(cfg *DataSplitConfig) (*DataSet, *DataSet, error) {
	var (
		trainingShare                float64
		method                       DataSplitMethod
		trainingRecords, testRecords []Record
	)

	if cfg == nil || cfg.TrainingShare == 0.0 {
		trainingShare = DEFAULT_TRAINING_SHARE
		method = SplitRandom
	} else {
		trainingShare = cfg.TrainingShare
		method = cfg.Method
	}

	splitPoint := int(float64(len(ds.Records)) * trainingShare)
	switch method {
	case SplitRandom:
		shuffled := randomShuffle(ds.Records)
		trainingRecords = shuffled[:splitPoint]
		testRecords = shuffled[splitPoint:]
	case SplitSequential:
		trainingRecords = ds.Records[:splitPoint]
		testRecords = ds.Records[splitPoint:]
	}

	training, _ := NewDataSet(ds.ClassNames, ds.AttributeNames, trainingRecords)
	test, _ := NewDataSet(ds.ClassNames, ds.AttributeNames, testRecords)
	return training, test, nil
}

func randomShuffle(source []Record) []Record {
	var deck []Record

	if len(source) > maxRecordCount {
		// Make a copy, we're not modifying the original record set
		deck = slices.Clone(source)
		// Split into subsets (cut the deck into smaller decks), shuffle, then merge (in reverse order)
		var subdecks [][]Record
		if len(deck)%maxRecordCount != 0 {
			subdecks = make([][]Record, (len(deck)/maxRecordCount)+1)
		} else {
			subdecks = make([][]Record, (len(deck) / maxRecordCount))
		}

		// They say a deck is randomized after seven shuffles...
		for p := 0; p < 7; p++ {
			if len(deck)%maxRecordCount == 0 {
				// Deck length is an even multiple of maxRecordCount, so when we "cut the cards",
				// we make the first (all but one) subdecks maxRecordCount in length, starting with index
				// maxRecordCount/2.
				for i := 0; i < len(subdecks)-1; i++ {
					subdecks[i] = shuffleDeck(deck[(i*maxRecordCount)+(maxRecordCount/2) : ((i+1)*maxRecordCount)+(maxRecordCount/2)])
				}

				// The last deck is the last maxRecordCount/2 cards + the first maxRecordCount/2 cards
				unsorted := append(deck[:maxRecordCount/2], deck[((len(subdecks)-1)*maxRecordCount)+(maxRecordCount/2):]...)
				subdecks[len(subdecks)-1] = shuffleDeck(unsorted)
			} else {
				for i, _ := range subdecks {
					if i != len(subdecks)-1 {
						subdecks[i] = shuffleDeck(deck[i*maxRecordCount : (i+1)*maxRecordCount])
					} else {
						subdecks[i] = shuffleDeck(deck[i*maxRecordCount:])
					}
				}
			}

			// range over all the subdecks but the last one.  We reverse the order in the merge
			// to make sure are not just recreating the subdecks with the same distributions.
			for i := 0; i < len(subdecks)-1; i++ {
				for j := range subdecks[i] {
					deck[(len(deck)-1)-((i*maxRecordCount)+j)] = subdecks[i][j]
				}
			}

			for j := range subdecks[len(subdecks)-1] {
				deck[len(subdecks[len(subdecks)-1])-(j+1)] = subdecks[len(subdecks)-1][j]
			}
		}
	} else {
		deck = shuffleDeck(source)
	}

	return deck
}

func shuffleDeck(source []Record) []Record {
	maxIdx := int(math.Pow(float64(len(source)), 3))

	randomizedSparse := make([]*Record, maxIdx+1)

	// Generate and store the list of randomized record indices
	sortKeys := make([]int, len(source))
	for i, _ := range sortKeys {
		// This is crappy and could bury us, but I don't feel like implementing wraparound logic
		// right now (incrementing the random index by 1 until we either find a free slot or wrap
		// around back to the beginning of the slice, then look for free slots from there)
		// This could theoretically take n (or more) iterations to find one free slot,
		// which is O(n^2)...or worse...over the set.  It's worse b/c there is no guarantee that
		// we don't retry the same slot (we don't exclude a slot we've checked and found occupied
		// from the random generation).  So...yeah, this needs to be fixed.
		for {
			rdIdx := rand.Intn(maxIdx + 1)
			if randomizedSparse[rdIdx] == nil {
				randomizedSparse[rdIdx] = &(source[i])
				sortKeys[i] = rdIdx
				break
			}
		}
	}

	// Condense to get the result
	randomized := make([]Record, len(source))
	slices.Sort(sortKeys)
	for i, k := range sortKeys {
		randomized[i] = *(randomizedSparse[k])
	}

	return randomized
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
