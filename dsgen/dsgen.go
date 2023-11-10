package dsgen

import (
	"errors"
	"math/rand"

	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/wyvern"
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

// ClassNames returns a slice containing the names of all classes defined
// in the configuration.
func (datasetConfig *DatasetConfig) ClassNames() []string {
	names := make([]string, 0)

	for name, _ := range datasetConfig.Classes {
		names = append(names, name)
	}

	return names
}

// Generate generates a new dataset from the specified configuration.
func GenerateDataset(datasetConfig *DatasetConfig) (*classifiers.DataSet, error) {
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
	// We could do this in batches if the count is excessive...
	records := make([]classifiers.Record, datasetConfig.RecordCount)
	for ri := range records {
		records[ri] = classifiers.Record{
			AttributeValues: make(wyvern.Vector[float64], len(attributeNames)),
		}
	}

	// Assign each record index a random class (note: this does NOT imply an even distribution among the classes)
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
			countsByQuintile, _ := ComputeQuintileDistribution(len(recordsInClass), attr.AllocationsByQuintile)
			recordsByQuintile, _ := AssignQuintiles(recordsInClass, countsByQuintile)
			for qi, quintile := range recordsByQuintile {
				for _, r := range quintile {
					records[r].AttributeValues[attrIdx] = ComputeAttributeValue(attr.LowerBound, attr.UpperBound, qi)
				}
			}
		}
	}

	return classifiers.NewDataSet(classNames, attributeNames, records)
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

// AssignQuintiles - Given a set of record indices ([]int), assign them to quintiles according to quintileCounts
func AssignQuintiles(indices []int, quintileCounts []int) ([][]int, error) {
	if len(quintileCounts) > 5 {
		return nil, errors.New("Too many quintiles (only 5 permitted) specified")
	}

	quintileRecords := make([][]int, 5)
	for i := range quintileRecords {
		quintileRecords[i] = make([]int, 0)
	}

	bucket := 0
	for ri := 0; ri < len(indices); ri++ {
		// Yeah, risky
		// Find the first quintile which is not at capacity, starting with the one
		// to which we most recently added.
		for b := bucket; ; b++ {
			quintileIdx := b % 5
			if len(quintileRecords[quintileIdx]) < quintileCounts[quintileIdx] {
				quintileRecords[quintileIdx] = append(quintileRecords[quintileIdx], indices[ri])
				bucket = b
				break
			}
		}
	}

	return quintileRecords, nil
}

// ComputeQuintileDistribution divides a set of record indices according to the
// given allocations (which, added together, must not exceed 100)
func ComputeQuintileDistribution(recordCount int, allocations []float64) ([]int, error) {
	if recordCount < 1 {
		return nil, errors.New("Record count must be > 0")
	}

	if len(allocations) > 5 {
		return nil, errors.New("Only five (5) allocations can be specified")
	}

	allocTotal := 0.0
	for _, a := range allocations {
		allocTotal += a
	}
	if allocTotal > 100.0 {
		return nil, errors.New("Allocations cannot total greater than 100.0")
	}

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

	return recordCountsByQuintile, nil
}

// ComputeAttributeValue computes a random value within the specified quintile.
// The first quintile (index 0) ranges from 0 to 20%.
// The other quntiles range from 21 to 40, 41 to 60, 61 to 80, and 81 to 100.
func ComputeAttributeValue(lower, upper float64, quintile int) float64 {
	var (
		percentage, rangeMagnitude float64
	)

	rangeMagnitude = upper - lower

	if quintile == 0 {
		percentage = float64((quintile*20)+rand.Intn(21)) / 100.0
	} else {
		percentage = float64((quintile*20)+rand.Intn(20)+1) / 100.0
	}

	return (rangeMagnitude * percentage) + lower
}
