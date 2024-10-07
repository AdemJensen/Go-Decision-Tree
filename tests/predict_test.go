package tests

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"DecisionTree/dataset"
	"DecisionTree/tree"
	"fmt"
	"log"
	"testing"
)

func TestPredict(t *testing.T) {
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
	dataset.PreProcessData(trainData)

	// read test data
	testData, err := data.ReadValues(config.Conf, attrTable, testDataFile)
	if err != nil {
		log.Fatalf("failed to read testing data: %v", err)
		return
	}
	dataset.PreProcessData(testData)
	print("OK\n")

	// read tree from file
	print("Reading tree...")
	tr, err := tree.ReadTreeFromFile("../tree.json")
	if err != nil {
		log.Fatalf("failed to read tree: %v", err)
		return
	}
	print("OK\n")

	fmt.Printf("=========================== TRAIN DATASET ===========================\n")

	// predict all test data, calculate accuracy
	res, err := tree.TestRun(tr, trainData)
	if err != nil {
		t.Fatalf("failed to do test run: %v", err)
		return
	}
	outputTestResult(tr, res)

	fmt.Printf("=========================== TEST DATASET ===========================\n")

	// predict all test data, calculate accuracy
	res, err = tree.TestRun(tr, testData)
	if err != nil {
		t.Fatalf("failed to do test run: %v", err)
		return
	}
	outputTestResult(tr, res)
}

func outputTestResult(tr *tree.Tree, res *tree.TestResults) {
	if tr != nil {
		fmt.Printf("Nodes count: %d\n", tr.GetNodeCount())
		fmt.Printf("Leaf Nodes count: %d\n", len(tr.GetLeafNodes()))
		fmt.Printf("Max depth: %d\n", tr.GetMaxDepth())
	}
	fmt.Printf("Accuracy: %.2f%%\n", res.Accuracy*100)
	fmt.Printf("Avg predict time: %s\n", res.AvgPredictTime.String())
	fmt.Printf("Pessimistic error: %.2f%%\n", res.PessimisticError*100)
	for class, count := range res.ClassDataCount {
		fmt.Printf("Class [%s] data frequency: %.2f%%\n", class, float64(count)/float64(res.TotalDataCount)*100)
		fmt.Printf("Class [%s] recall: %.2f%%\n", class, res.ClassRecall[class]*100)
		fmt.Printf("Class [%s] precision: %.2f%%\n", class, res.ClassPrecision[class]*100)
	}
	fmt.Printf("Confusion matrix:\n")
	for actual, row := range res.ConfusionMatrix {
		for predicted, count := range row {
			fmt.Printf("Actual [%s] & Predict [%s]: %d\n", actual, predicted, count)
		}
	}
}
