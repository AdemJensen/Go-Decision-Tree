package tree

import (
	"DecisionTree/config"
	"sort"
)

func splitInstancesByContinuousAttr(_ *config.Config, rootEntropy float64, attrIndex int, instances []*WeightedInstance) ([]*Node, float64, error) {

	// Calculate count for all class values
	var (
		nonMissingInstances     []*WeightedInstance
		missingInstances        []*WeightedInstance
		classValueCount         = make(map[string]float64) // will only count non-missing instances
		nonMissingInstanceCount = 0.0
		missingInstanceCount    = 0.0
	)
	for _, instance := range instances {
		if instance.Instance.AttributeValues[attrIndex].IsMissing() {
			missingInstances = append(missingInstances, instance)
			missingInstanceCount += instance.Weight
			continue
		}
		nonMissingInstances = append(nonMissingInstances, instance)
		classValue := instance.Instance.ClassValue.Value().(string)
		classValueCount[classValue] += instance.Weight
		nonMissingInstanceCount += instance.Weight
	}

	if len(nonMissingInstances) < 2 {
		return nil, 0, nil
	}

	// sort instances for continuous attribute
	sort.Slice(nonMissingInstances, func(i, j int) bool {
		return nonMissingInstances[i].Instance.AttributeValues[attrIndex].Value().(float64) < nonMissingInstances[j].Instance.AttributeValues[attrIndex].Value().(float64)
	})

	// from left to right, calculate the best split
	var (
		bestSplitValue     float64
		bestSplitGain      float64
		bestSplitPoint     = 0
		leftClassCnt       = make(map[string]float64)
		rightClassCnt      = classValueCount
		leftInstanceCount  float64
		rightInstanceCount = nonMissingInstanceCount
	)
	for i := 1; i < len(instances); i++ {
		// update class count
		leftClassCnt[instances[i-1].Instance.ClassValue.Value().(string)] += instances[i-1].Weight
		rightClassCnt[instances[i-1].Instance.ClassValue.Value().(string)] -= instances[i-1].Weight
		leftInstanceCount += instances[i-1].Weight
		rightInstanceCount -= instances[i-1].Weight

		v1 := instances[i-1].Instance.AttributeValues[attrIndex].Value().(float64)
		v2 := instances[i].Instance.AttributeValues[attrIndex].Value().(float64)
		if v1 == v2 {
			continue
		}

		// calculate gain for split
		entropyLeft := calculateEntropy(instances[:i], leftInstanceCount, leftClassCnt)
		entropyRight := calculateEntropy(instances[i:], rightInstanceCount, rightClassCnt)
		entropy := (entropyLeft*float64(i) + entropyRight*float64(len(instances)-i)) / float64(len(instances))

		gain := (rootEntropy - entropy) * float64(len(nonMissingInstances)) / float64(len(instances))
		if gain > bestSplitGain {
			bestSplitGain = gain
			bestSplitValue = (v1 + v2) / 2
			bestSplitPoint = i
		}
	}

	if bestSplitPoint == 0 {
		return nil, 0, nil
	}

	// split instances
	leftInstances := nonMissingInstances[:bestSplitPoint]
	rightInstances := nonMissingInstances[bestSplitPoint:]
	// distribute missing value instances
	for _, instance := range missingInstances {
		newLeftInstance := instance.CopyWithScale(leftInstanceCount / nonMissingInstanceCount)
		newRightInstance := instance.CopyWithScale(rightInstanceCount / nonMissingInstanceCount)
		leftInstances = append(leftInstances, newLeftInstance)
		rightInstances = append(rightInstances, newRightInstance)
	}
	return []*Node{
		{
			Condition:     newLessThanCondition(instances[0].Instance.AttributeValues[attrIndex].Attribute(), bestSplitValue),
			instances:     leftInstances,
			IsPrioritized: len(leftInstances) >= len(rightInstances),
		},
		{
			Condition:     newGreaterThanEqCondition(instances[0].Instance.AttributeValues[attrIndex].Attribute(), bestSplitValue),
			instances:     rightInstances,
			IsPrioritized: len(rightInstances) > len(leftInstances),
		},
	}, bestSplitGain, nil
}
