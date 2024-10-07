package main

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"DecisionTree/tree"
	"DecisionTree/utils"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/gosuri/uiprogress"
)

func TestHyperParams(t *testing.T) {
	var (
		attributesFile = "dataset/adult.names"
		trainDataFile  = "dataset/adult.data"
		testDataFile   = "dataset/adult.test"

		maxDepthList                                = []int{30, 40, 50, 60, 70}
		minSamplesSplitList                         = []int{8, 16, 32}
		minSamplesLeafList                          = []int{4, 8, 16}
		minImpurityDecreaseList                     = []float64{0.01, 0.05, 0.1, 0.2}
		maxNominalBruteForceScaleList               = []int{8, 12, 16}
		minPostPruneGeneralizationErrorDecreaseList = []float64{0.0}
	)
	uiProgress := config.GetUiProgress()

	// load dataset
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
	preProcessData(trainData)

	// read test data
	testData, err := data.ReadValues(config.Conf, attrTable, testDataFile)
	if err != nil {
		log.Fatalf("failed to read testing data: %v", err)
		return
	}
	preProcessData(testData)
	print("OK\n")

	var (
		bar = uiProgress.AddBar(len(maxDepthList) * len(minSamplesSplitList) * len(minSamplesLeafList) * len(minImpurityDecreaseList) * len(maxNominalBruteForceScaleList) * len(minPostPruneGeneralizationErrorDecreaseList)).PrependFunc(func(b *uiprogress.Bar) string {
			return "Performing Tests"
		}).AppendFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("%d/%d, %.2f%%, Time Elapsed: %s", b.Current(), b.Total, b.CompletedPercent(), b.TimeElapsedString())
		})
		allResults []*HyperParamTestResult
	)

	for _, maxDepth := range maxDepthList {
		for _, minSamplesSplit := range minSamplesSplitList {
			for _, minSamplesLeaf := range minSamplesLeafList {
				for _, minImpurityDecrease := range minImpurityDecreaseList {
					for _, maxNominalBruteForceScale := range maxNominalBruteForceScaleList {
						for _, minPostPruneGeneralizationErrorDecrease := range minPostPruneGeneralizationErrorDecreaseList {
							bar.Incr()
							uiProgress.Bars = uiProgress.Bars[:1]
							conf := &config.Config{
								ConsiderInvalidDataAsMissing:            true,
								MaxDepth:                                maxDepth,
								MinSamplesSplit:                         minSamplesSplit,
								MinSamplesLeaf:                          minSamplesLeaf,
								MinImpurityDecrease:                     minImpurityDecrease,
								MaxNominalBruteForceScale:               maxNominalBruteForceScale,
								MinPostPruneGeneralizationErrorDecrease: minPostPruneGeneralizationErrorDecrease,
								VerboseLog:                              false,
								LogFile:                                 "",
							}
							start := time.Now()
							tr, err := tree.BuildTree(conf, trainData)
							if err != nil {
								log.Printf("failed to build tree on conf %s: %v", utils.Json(conf), err)
								continue
							}
							trainTime := time.Since(start)
							res, err := tree.TestRun(tr, testData)
							if err != nil {
								log.Printf("failed to test tree %s: %v", utils.Json(conf), err)
								continue
							}
							allResults = append(allResults, &HyperParamTestResult{
								Conf:        conf,
								TrainTime:   trainTime,
								NNodes:      tr.GetNodeCount(),
								NLeaf:       len(tr.GetLeafNodes()),
								TestMetrics: res,
							})
						}
					}
				}
			}
		}
	}

	// sort and get the best result
	slices.SortFunc(allResults, func(a, b *HyperParamTestResult) int {
		switch {
		case a.TestMetrics.Accuracy > b.TestMetrics.Accuracy:
			return -1
		case a.TestMetrics.Accuracy < b.TestMetrics.Accuracy:
			return 1
		default:
			return 0
		}
	})
	bestResult := allResults[0]
	fmt.Printf("Best Config: %s\n", utils.JsonPretty(bestResult.Conf))
	fmt.Printf("Best Result:\n")
	outputTestResult(bestResult.TestMetrics)

	// save all test content to file
	file, err := os.Create("hyper_param_test.json")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, err = file.WriteString(utils.JsonPretty(allResults))
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
		return
	}
}

type HyperParamTestResult struct {
	Conf        *config.Config
	TrainTime   time.Duration
	NNodes      int
	NLeaf       int
	TestMetrics *tree.TestResults
}
