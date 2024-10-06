package tree

import (
	"DecisionTree/data"
	"math"
)

// calculateEntropy calculates the entropy of a set of instances.
// The fastClassFrequencyCnt is the frequency count of the majority class in the parent node, input a non-empty list set a fast-forward frequency.
func calculateEntropy(instances []*WeightedInstance, allInstanceCount float64, fastClassFrequencyCnt map[string]float64) float64 {
	if len(instances) == 0 {
		return 0
	}

	// if fastClassFrequency is empty, calculate frequency
	if fastClassFrequencyCnt == nil {
		fastClassFrequencyCnt = calculateClassFrequencyCnt(instances)
	}

	if allInstanceCount == 0 {
		for _, instance := range instances {
			allInstanceCount += instance.Weight
		}
	}

	// calculate entropy
	classValues := instances[0].Instance.ClassValue.Attribute().(*data.NominalAttribute).AcceptedValues
	entropy := 0.0
	for _, value := range classValues {
		frequency := fastClassFrequencyCnt[value] / allInstanceCount
		if frequency == 0 {
			continue
		}
		entropy -= frequency * math.Log2(frequency)
	}
	return entropy
}

func calculateClassFrequencyCnt(instances []*WeightedInstance) map[string]float64 {
	// calculate frequency
	valueCount := make(map[string]float64)
	totalCount := 0.0
	for _, instance := range instances {
		valueCount[instance.Instance.ClassValue.Value().(string)] += instance.Weight
		totalCount += instance.Weight
	}
	return valueCount
}
