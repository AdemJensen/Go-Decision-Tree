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
		// testDataFile   = "dataset/adult.test"
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

	// pre-process data
	preProcessData(trainData)

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
}

func preProcessData(valueTable *data.ValueTable) {
	// 75% is <=50K, 25% is >50K
	// Resample the data to make it balanced
	dataInstances := make(map[string][]*data.Instance)
	for _, attr := range valueTable.Instances {
		classValue := attr.ClassValue.Value().(string)
		dataInstances[classValue] = append(dataInstances[classValue], attr)
	}
	valueTable.Instances = append(valueTable.Instances, dataInstances[">50K"]...)
	valueTable.Instances = append(valueTable.Instances, dataInstances[">50K"]...)
}
