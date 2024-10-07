package tests

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"DecisionTree/dataset"
	"DecisionTree/tree"
	"DecisionTree/utils"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/gosuri/uiprogress"
)

func TestHyperParams(t *testing.T) {
	var (
		attributesFile = "../dataset/adult.names"
		trainDataFile  = "../dataset/adult.data"
		testDataFile   = "../dataset/adult.test"

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
	dataset.PreProcessData(trainData)

	// read test data
	testData, err := data.ReadValues(config.Conf, attrTable, testDataFile)
	if err != nil {
		log.Fatalf("failed to read testing data: %v", err)
		return
	}
	dataset.PreProcessData(testData)
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
	outputTestResult(nil, bestResult.TestMetrics)

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

func aggregateHyperParamTestResults(results []*HyperParamTestResult) *HyperParamTestResult {
	if len(results) == 0 {
		return nil
	}
	totalTrainTime := time.Duration(0)
	totalNNodes := 0
	totalNLeaf := 0
	totalTestMetrics := &tree.TestResults{
		Accuracy:       0,
		ClassRecall:    make(map[string]float64),
		ClassPrecision: make(map[string]float64),
		AvgPredictTime: 0,
	}

	for _, result := range results {
		totalTrainTime += result.TrainTime
		totalNNodes += result.NNodes
		totalNLeaf += result.NLeaf
		totalTestMetrics.Accuracy += result.TestMetrics.Accuracy
		totalTestMetrics.AvgPredictTime += result.TestMetrics.AvgPredictTime
		for k, v := range result.TestMetrics.ClassRecall {
			totalTestMetrics.ClassRecall[k] += v
		}
		for k, v := range result.TestMetrics.ClassPrecision {
			totalTestMetrics.ClassPrecision[k] += v
		}
	}

	// calculate the average
	return &HyperParamTestResult{
		Conf:      results[0].Conf,
		TrainTime: totalTrainTime / time.Duration(len(results)),
		NNodes:    totalNNodes / len(results),
		NLeaf:     totalNLeaf / len(results),
		TestMetrics: &tree.TestResults{
			Accuracy: totalTestMetrics.Accuracy / float64(len(results)),
			ClassRecall: func() map[string]float64 {
				m := make(map[string]float64)
				for k, v := range totalTestMetrics.ClassRecall {
					m[k] = v / float64(len(results))
				}
				return m
			}(),
			ClassPrecision: func() map[string]float64 {
				m := make(map[string]float64)
				for k, v := range totalTestMetrics.ClassPrecision {
					m[k] = v / float64(len(results))
				}
				return m
			}(),
			AvgPredictTime: totalTestMetrics.AvgPredictTime / time.Duration(len(results)),
		},
	}
}

func TestPlotHyperParamsGraph(t *testing.T) {
	// read json content from file
	file, err := os.Open("hyper_param_test.json")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
		return
	}

	// 反序列化 JSON 内容
	var results []*HyperParamTestResult
	if err := json.Unmarshal(bytes, &results); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
		return
	}

	renderLineChart("max_depth", results, func(result *HyperParamTestResult) float64 {
		return float64(result.Conf.MaxDepth)
	})
	renderLineChart("min_samples_split", results, func(result *HyperParamTestResult) float64 {
		return float64(result.Conf.MinSamplesSplit)
	})
	renderLineChart("min_impurity_decrease", results, func(result *HyperParamTestResult) float64 {
		return result.Conf.MinImpurityDecrease
	})
	renderLineChart("min_samples_leaf", results, func(result *HyperParamTestResult) float64 {
		return float64(result.Conf.MinSamplesLeaf)
	})
}

func renderLineChart(name string, results []*HyperParamTestResult, xAxis func(result *HyperParamTestResult) float64) {
	// Create a line chart
	line := charts.NewLine()

	// aggregate the data
	var aggData = make(map[float64][]*HyperParamTestResult)
	for _, result := range results {
		x := xAxis(result)
		aggData[x] = append(aggData[x], result)
	}
	results = nil
	for _, res := range aggData {
		results = append(results, aggregateHyperParamTestResults(res))
	}
	slices.SortFunc(results, func(a, b *HyperParamTestResult) int {
		va, vb := xAxis(a), xAxis(b)
		return cmp.Compare[float64](va, vb)
	})

	// X-axis: MaxDepth values
	var xValues []float64
	for _, result := range results {
		xValues = append(xValues, xAxis(result))
	}
	line.SetXAxis(xValues)

	// Y-axis: Different metrics
	var accuracies []opts.LineData
	// var trainTimes []opts.LineData
	var recallLessThan50K, recallGreaterThan50K []opts.LineData
	var precisionLessThan50K, precisionGreaterThan50K []opts.LineData

	for _, result := range results {
		accuracies = append(accuracies, opts.LineData{Value: result.TestMetrics.Accuracy})
		//trainTimes = append(trainTimes, opts.LineData{Value: result.TrainTime.Seconds()})
		recallLessThan50K = append(recallLessThan50K, opts.LineData{Value: result.TestMetrics.ClassRecall["<=50K"]})
		recallGreaterThan50K = append(recallGreaterThan50K, opts.LineData{Value: result.TestMetrics.ClassRecall[">50K"]})
		precisionLessThan50K = append(precisionLessThan50K, opts.LineData{Value: result.TestMetrics.ClassPrecision["<=50K"]})
		precisionGreaterThan50K = append(precisionGreaterThan50K, opts.LineData{Value: result.TestMetrics.ClassPrecision[">50K"]})
	}

	// Add the data to the chart
	True := true
	line.AddSeries("Accuracy", accuracies).
		//AddSeries("TrainTime (seconds)", trainTimes).
		AddSeries("Recall <=50K", recallLessThan50K).
		AddSeries("Precision <=50K", precisionLessThan50K).
		AddSeries("Recall >50K", recallGreaterThan50K).
		AddSeries("Precision >50K", precisionGreaterThan50K).
		SetGlobalOptions(
			charts.WithLegendOpts(opts.Legend{
				Show: &True,
			}),
			charts.WithTooltipOpts(opts.Tooltip{
				Show: &True,
			}),
		)

	// Render the chart to an HTML file
	f, err := os.Create(fmt.Sprintf("../docs/hyper_param_results_%s.html", name))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_ = line.Render(f)
}
