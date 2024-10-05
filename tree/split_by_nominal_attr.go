package tree

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"fmt"
	"math"
	"slices"
)

func splitInstancesByNominalAttr(conf *config.Config, rootEntropy float64, attrIndex int, instances []*WeightedInstance) ([]*nominalSplitUnit, float64, error) {
	if len(instances) == 0 {
		return nil, 0, nil
	}

	var (
		bestGain  = 0.0
		bestSplit []*nominalSplitUnit
	)

	// try multi-way split first
	multiWaySplit, multiWayGain, err := multiWaySplitByNominalAttr(conf, rootEntropy, attrIndex, instances)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to do multi-way split: %w", err)
	}
	if multiWaySplit != nil && multiWayGain > bestGain {
		bestGain = multiWayGain
		bestSplit = multiWaySplit
	}

	// try binary split
	binarySplit, binaryGain, err := binarySplitByNominalAttr(conf, rootEntropy, attrIndex, instances)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to do binary split: %w", err)
	}
	if binarySplit != nil && binaryGain > bestGain {
		bestGain = binaryGain
		bestSplit = binarySplit
	}

	return bestSplit, bestGain, nil
}

func multiWaySplitByNominalAttr(conf *config.Config, rootEntropy float64, attrIndex int, instances []*WeightedInstance) ([]*nominalSplitUnit, float64, error) {
	classifiedInstancesMap, missingValueInstances := classifyInstancesByNominalAttr(instances, attrIndex)
	var classifyUnits []*nominalSplitUnit
	for value, instanceList := range classifiedInstancesMap {
		classifyUnits = append(classifyUnits, newNominalValueUnit(value, instanceList))
	}
	// calculate gain for this split
	gain := calculateGainForNominalSplit(rootEntropy, instances, classifyUnits)
	// distribute missing value instances
	distributeMissingValuesToNominalSplit(classifyUnits, missingValueInstances)
	if !checkNominalSplitMinSamplesLeaf(conf, classifyUnits) {
		return nil, 0, nil
	}
	return classifyUnits, gain, nil
}

// binarySplitByNominalAttr splits instances into two groups by a nominal attribute.
// For nominal attribute, if the number of accepted values is less than max_nominal_brute_force_scale, use
// brute-force to find the best split. If not, we will first join the values with fewer instances until the
// number of values is less than or equal to max_nominal_brute_force_scale, then perform brute-force.
// returns: split result, gain, error
func binarySplitByNominalAttr(conf *config.Config, rootEntropy float64, attrIndex int, instances []*WeightedInstance) ([]*nominalSplitUnit, float64, error) {
	classifiedInstancesMap, missingValueInstances := classifyInstancesByNominalAttr(instances, attrIndex)
	// join values with fewer instances until the number of values is less than or equal to max_nominal_brute_force_scale
	// initialize a join list
	var sortUnits []*nominalSplitUnit
	for value, instanceList := range classifiedInstancesMap {
		sortUnits = append(sortUnits, newNominalValueUnit(value, instanceList))
	}
	// sort the join list by the number of instances
	sortNominalValueUnitList(sortUnits)
	for len(sortUnits) > conf.MaxNominalBruteForceScale {
		sortUnits[0] = joinNominalValueUnit(sortUnits[0], sortUnits[1])
		sortUnits = append(sortUnits[:1], sortUnits[2:]...)
		// sort after joining
		sortNominalValueUnitList(sortUnits)
	}

	///////////////////////////////////////////////////////////////////////
	// perform brute-force, find the best split
	///////////////////////////////////////////////////////////////////////
	totalSplits := int(math.Pow(2, float64(len(sortUnits)))) - 1
	var (
		bestGain  = 0.0
		bestSplit []*nominalSplitUnit
	)

	// Iterate over all possible splits (represented as bitmasks)
	for i := 1; i < totalSplits; i++ {
		var (
			left  *nominalSplitUnit
			right *nominalSplitUnit
		)

		// Use the bitmask to partition the elements into two subsets
		for j := 0; j < len(sortUnits); j++ {
			if i&(1<<j) != 0 {
				left = joinNominalValueUnit(left, sortUnits[j])
			} else {
				right = joinNominalValueUnit(right, sortUnits[j])
			}
		}

		// Ensure both subsets are non-empty
		if left == nil || right == nil {
			continue
		}

		splitRes := []*nominalSplitUnit{left, right}

		// if the number of instances in any subset is less than min_samples_leaf, skip this split
		if !checkNominalSplitMinSamplesLeaf(conf, splitRes) {
			continue
		}

		// Calculate entropy for each subset, if the sum of the entropy is less than the current best split, update the best split
		gain := calculateGainForNominalSplit(rootEntropy, instances, splitRes)
		if gain > bestGain {
			bestGain = gain
			bestSplit = []*nominalSplitUnit{left, right}
		}
	}
	if len(bestSplit) != 2 {
		return nil, 0, nil
	}
	// distribute missing value instances
	distributeMissingValuesToNominalSplit(bestSplit, missingValueInstances)
	return bestSplit, bestGain, nil
}

// classifyInstancesByNominalAttr classifies instances by a nominal attribute.
// Returns: a map of accepted values to instances, and instances with missing values.
func classifyInstancesByNominalAttr(instances []*WeightedInstance, attrIndex int) (map[string][]*WeightedInstance, []*WeightedInstance) {
	var (
		res              = make(map[string][]*WeightedInstance)
		missingInstances []*WeightedInstance
	)
	for _, instance := range instances {
		value := instance.Instance.AttributeValues[attrIndex]
		if value.IsMissing() {
			missingInstances = append(missingInstances, instance)
			continue
		}
		res[value.Value().(string)] = append(res[value.Value().(string)], instance)
	}
	return res, missingInstances
}

type nominalSplitUnit struct {
	values                  []string
	instances               []*WeightedInstance
	count                   float64 // count considers weight
	classValueInstanceCount map[string]float64
}

func newNominalValueUnit(value string, instances []*WeightedInstance) *nominalSplitUnit {
	res := &nominalSplitUnit{
		values:                  []string{value},
		instances:               instances,
		count:                   SumInstanceWeights(instances),
		classValueInstanceCount: make(map[string]float64),
	}
	for _, instance := range instances {
		res.classValueInstanceCount[instance.Instance.ClassValue.Value().(string)] += instance.Weight
	}
	return res
}

func joinNominalValueUnit(a, b *nominalSplitUnit) *nominalSplitUnit {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	res := &nominalSplitUnit{
		values:                  append(a.values, b.values...),
		instances:               append(a.instances, b.instances...),
		count:                   a.count + b.count,
		classValueInstanceCount: make(map[string]float64),
	}
	for k, v := range a.classValueInstanceCount {
		res.classValueInstanceCount[k] += v
	}
	for k, v := range b.classValueInstanceCount {
		res.classValueInstanceCount[k] += v
	}
	return res
}

func sortNominalValueUnitList(units []*nominalSplitUnit) {
	slices.SortFunc(units, func(a, b *nominalSplitUnit) int {
		switch {
		case a.count < b.count:
			return -1
		case a.count > b.count:
			return 1
		default:
			return 0
		}
	})
}

func calculateGainForNominalSplit(rootEntropy float64, instances []*WeightedInstance, split []*nominalSplitUnit) float64 {
	var (
		entropy    = 0.0
		totalCount = 0.0
	)
	for _, unit := range split {
		totalCount += unit.count
	}
	for _, unit := range split {
		entropy += calculateEntropy(unit.instances, unit.count, unit.classValueInstanceCount) * unit.count / totalCount
	}
	instanceWeightSum := SumInstanceWeights(instances)
	return (rootEntropy - entropy) * totalCount / instanceWeightSum
}

func distributeMissingValuesToNominalSplit(split []*nominalSplitUnit, missingValueInstances []*WeightedInstance) {
	// get the total count of non-missing value instances
	totalCount := 0.0
	for _, unit := range split {
		totalCount += unit.count
	}
	for _, instance := range missingValueInstances {
		for i, unit := range split {
			newInstance := instance.CopyWithScale(unit.count / totalCount)
			split[i] = joinNominalValueUnit(unit, newNominalValueUnit("", []*WeightedInstance{newInstance}))
		}
	}
}

func checkNominalSplitMinSamplesLeaf(conf *config.Config, split []*nominalSplitUnit) bool {
	for _, unit := range split {
		if len(unit.instances) < conf.MinSamplesLeaf {
			return false
		}
	}
	return true
}

func buildNodeListFromNominalSplit(attribute data.Attribute, split []*nominalSplitUnit) []*Node {
	var res []*Node
	for _, unit := range split {
		res = append(res, &Node{
			Condition: newIsOneOfCondition(attribute, unit.values),
			Children:  nil,
			instances: unit.instances,
		})
	}
	return res
}
