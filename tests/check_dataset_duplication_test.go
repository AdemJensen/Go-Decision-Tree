package tests

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"log"
	"testing"
)

func TestCheckDatasetDuplication(t *testing.T) {
	// There might be redundant data in the dataset.
	// Most significantly, education and education-num are probably the same.
	// We will use this program if that is the case.
	// If so, we will use education-num and remove education.

	var (
		attributesFile = "../dataset/adult.names"
		trainDataFile  = "../dataset/adult.data"
		testDataFile   = "../dataset/adult.test"
	)
	// read data attributes
	print("Reading dataset...")
	attrTable, err := data.ReadAttributes(attributesFile)
	if err != nil {
		log.Fatalf("failed to read attributes: %v", err)
		return
	}

	// read train data
	trainData, err := data.ReadValues(config.Conf, attrTable, trainDataFile)
	if err != nil {
		log.Fatalf("failed to read training data: %v", err)
		return
	}

	// read test data
	testData, err := data.ReadValues(config.Conf, attrTable, testDataFile)
	if err != nil {
		log.Fatalf("failed to read testing data: %v", err)
		return
	}
	print("OK\n")

	// Check the relationship between education and education-num
	eduMapping := make(map[string]float64)
	targetAttr1 := attrTable.GetAttrByName("education")
	targetAttr2 := attrTable.GetAttrByName("education-num")
	for _, instance := range append(trainData.Instances, testData.Instances...) {
		edu := instance.GetValueByAttr(targetAttr1)
		eduNum := instance.GetValueByAttr(targetAttr2)
		if edu.IsMissing() || eduNum.IsMissing() {
			continue
		}
		if _, ok := eduMapping[edu.Value().(string)]; ok {
			if eduMapping[edu.Value().(string)] != eduNum.Value().(float64) {
				t.Logf("education and education-num are not the same (%.2f, %.2f)", eduMapping[edu.Value().(string)], eduNum.Value().(float64))
				return
			}
		} else {
			eduMapping[edu.Value().(string)] = eduNum.Value().(float64)
		}
	}
	t.Logf("education and education-num are the same")
}
