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

	//_, err = tr.Predict(testData.Instances[89])
	//if err != nil {
	//	t.Fatalf("failed to predict: %v", err)
	//	return
	//}

	// predict all test data, calculate accuracy
	var (
		correctCount      int
		classDataCount    = make(map[string]int)
		classCorrectCount = make(map[string]int)
	)
	for _, instance := range testData.Instances {
		classDataCount[instance.ClassValue.Value().(string)]++
		res, err := tr.Predict(instance)
		if err != nil {
			t.Fatalf("failed to predict: %v", err)
			return
		}
		if res == instance.ClassValue.Value().(string) {
			correctCount++
			classCorrectCount[res]++
		}
		//t.Logf("Completed predict for %d/%d instances", i+1, len(testData.Instances))
	}
	accuracy := float64(correctCount) / float64(len(testData.Instances))
	t.Logf("Accuracy: %.2f%%", accuracy*100)
	for class, count := range classDataCount {
		t.Logf("Class Data [%s]: %.2f%%", class, float64(count)/float64(len(testData.Instances))*100)
		t.Logf("Within Data [%s] Accuricy: %.2f%%", class, float64(classCorrectCount[class])/float64(count)*100)
	}
}
