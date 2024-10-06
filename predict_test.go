package main

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"DecisionTree/tree"
	"log"
	"testing"
)

func TestPredict(t *testing.T) {
	var (
		attributesFile = "dataset/adult.names"
		// trainDataFile  = "dataset/adult.data"
		testDataFile = "dataset/adult.test"
	)
	// read data attributes
	print("Reading dataset...")
	attrTable, err := data.ReadAttributes(attributesFile)
	if err != nil {
		log.Fatalf("failed to read attributes: %v", err)
		return
	}

	// read test data
	testData, err := data.ReadValues(config.Conf, attrTable, testDataFile)
	if err != nil {
		log.Fatalf("failed to read training data: %v", err)
		return
	}
	print("OK\n")

	// read tree from file
	print("Reading tree...")
	tr, err := tree.ReadTreeFromFile("tree.json")
	if err != nil {
		log.Fatalf("failed to read tree: %v", err)
		return
	}
	print("OK\n")

	// predict all test data, calculate accuracy
	res, err := tree.TestRun(tr, testData)
	if err != nil {
		t.Fatalf("failed to do test run: %v", err)
		return
	}
	t.Logf("Accuracy: %.2f%%", res.Accuracy*100)
	t.Logf("Pessimistic error: %.2f%%", res.PessimisticError*100)
	for class, count := range res.ClassDataCount {
		t.Logf("Class Data [%s] frequency: %.2f%%", class, float64(count)/float64(len(testData.Instances))*100)
		t.Logf("Within class [%s] predict accuricy: %.2f%%", class, float64(res.ClassCorrectCount[class])/float64(count)*100)
	}
}
