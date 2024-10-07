package tree

import (
	"DecisionTree/data"
	"fmt"
	"time"
)

type TestResults struct {
	TotalDataCount    int
	CorrectCount      int
	ErrorCount        int
	Accuracy          float64
	ClassDataCount    map[string]int
	ClassPredictCount map[string]int
	ClassCorrectCount map[string]int
	ClassErrorCount   map[string]int
	ClassRecall       map[string]float64
	ClassPrecision    map[string]float64
	ConfusionMatrix   map[string]map[string]int // actual class -> predicted class -> count
	PessimisticError  float64
	AvgPredictTime    time.Duration
}

func TestRun(tr *Tree, dataTable *data.ValueTable) (*TestResults, error) {
	return testRunNode(tr.RootNode, dataTable.Instances)
}

func testRunNode(node *Node, instances []*data.Instance) (*TestResults, error) {
	var (
		correctCount      int
		errorCount        int
		classDataCount    = make(map[string]int)
		classPredictCount = make(map[string]int)
		classCorrectCount = make(map[string]int)
		classErrorCount   = make(map[string]int)
		classRecall       = make(map[string]float64)
		classPrecision    = make(map[string]float64)
		confusionMatrix   = make(map[string]map[string]int)
		startTime         time.Time
	)
	startTime = time.Now()
	for i, instance := range instances {
		classDataCount[instance.ClassValue.Value().(string)]++
		res, err := node.Predict(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to predict instance %d: %w", i, err)
		}
		classPredictCount[res]++
		if _, ok := confusionMatrix[instance.ClassValue.Value().(string)]; !ok {
			confusionMatrix[instance.ClassValue.Value().(string)] = make(map[string]int)
		}
		confusionMatrix[instance.ClassValue.Value().(string)][res]++
		if res == instance.ClassValue.Value().(string) {
			correctCount++
			classCorrectCount[instance.ClassValue.Value().(string)]++
		} else {
			errorCount++
			classErrorCount[instance.ClassValue.Value().(string)]++
		}
	}
	accuracy := float64(correctCount) / float64(len(instances))
	for k, v := range classDataCount {
		classRecall[k] = float64(classCorrectCount[k]) / float64(v)
		classPrecision[k] = float64(classCorrectCount[k]) / float64(classPredictCount[k])
	}
	leafNodes := node.GetLeafNodes()
	var avgPredictTime time.Duration
	if len(instances) > 0 {
		avgPredictTime = time.Since(startTime) / time.Duration(len(instances))
	}
	return &TestResults{
		TotalDataCount:    len(instances),
		CorrectCount:      correctCount,
		ErrorCount:        errorCount,
		Accuracy:          accuracy,
		ClassDataCount:    classDataCount,
		ClassPredictCount: classPredictCount,
		ClassCorrectCount: classCorrectCount,
		ClassErrorCount:   classErrorCount,
		ClassRecall:       classRecall,
		ClassPrecision:    classPrecision,
		ConfusionMatrix:   confusionMatrix,
		PessimisticError:  calculatePessimisticError(errorCount, len(leafNodes), len(instances)),
		AvgPredictTime:    avgPredictTime,
	}, nil
}

func calculatePessimisticError(errorCount, leafNodesCount, totalDataCount int) float64 {
	return (float64(errorCount) + float64(leafNodesCount)*0.5) / float64(totalDataCount)
}
