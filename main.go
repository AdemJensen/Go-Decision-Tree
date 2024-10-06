package main

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"DecisionTree/tree"
	"log"
)

func main() {
	var (
		attributesFile = "dataset/adult.names"
		trainDataFile  = "dataset/adult.data"
		// testDataFile   = "data/adult.test"
	)
	// read data attributes
	print("Reading dataset...")
	attrTable, err := data.ReadAttributes(attributesFile)
	if err != nil {
		log.Fatalf("failed to read attributes: %v", err)
		return
	}

	// read training data
	trainData, err := data.ReadValues(config.Conf, attrTable, trainDataFile)
	if err != nil {
		log.Fatalf("failed to read training data: %v", err)
		return
	}
	print("OK\n")

	// train decision tree
	print("Training decision tree...")
	t, err := tree.BuildTree(config.Conf, trainData)
	if err != nil {
		log.Fatalf("failed to build tree: %v", err)
		return
	}
	print("OK\n")

	// save tree to file
	print("Saving tree...")
	err = tree.WriteTreeToFile(t, "tree.json")
	if err != nil {
		log.Fatalf("failed to save tree: %v", err)
		return
	}
	print("OK\n")

	// predict
	res, err := t.Predict(trainData.Instances[0])
	if err != nil {
		log.Fatalf("failed to predict: %v", err)
		return
	}
	println("Predict result: " + res)
}
