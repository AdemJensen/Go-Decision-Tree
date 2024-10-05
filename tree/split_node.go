package tree

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"fmt"
)

func splitNode(conf *config.Config, level int, node *Node) error {
	// if reach max depth, stop split
	if level >= conf.MaxDepth {
		return nil
	}

	// if reach min samples split, stop split
	if len(node.instances) < conf.MinSamplesSplit {
		return nil
	}

	// if all instances have the same class value, stop split
	if allSameClassValue(node.instances) {
		return nil
	}

	// for each attribute, try to split node, find the best split
	// Initial best split: do not split. Impurity is the impurity of the node
	var (
		bestSplitChildren []*Node // empty node list, means do not split
		bestSplitGain     = 0.0
		nodeEntropy       = calculateEntropy(node.instances, 0, nil)
	)
	for i, value := range node.instances[0].Instance.AttributeValues {
		attribute := value.Attribute()
		switch attribute.Type() {
		case data.Continuous:
			bestContinuousSplit, bestContinuousGain, err := splitInstancesByContinuousAttr(conf, nodeEntropy, i, node.instances)
			if err != nil {
				return fmt.Errorf("failed to split instances by continuous attribute: %w", err)
			}
			if len(bestContinuousSplit) > 0 && bestContinuousGain > bestSplitGain {
				bestSplitGain = bestContinuousGain
				bestSplitChildren = bestContinuousSplit
			}
		case data.Nominal:
			bestNominalSplit, bestNominalGain, err := splitInstancesByNominalAttr(conf, nodeEntropy, i, node.instances)
			if err != nil {
				return fmt.Errorf("failed to split instances by nominal attribute: %w", err)
			}
			if len(bestNominalSplit) > 0 && bestNominalGain > bestSplitGain {
				bestSplitGain = bestNominalGain
				bestSplitChildren = buildNodeListFromNominalSplit(attribute, bestNominalSplit)
			}
		}
	}

	// if no split, stop split
	if len(bestSplitChildren) == 0 {
		return nil
	}

	// split node
	node.Children = bestSplitChildren
	for _, child := range bestSplitChildren {
		// recursively split child node
		if err := splitNode(conf, level+1, child); err != nil {
			return fmt.Errorf("failed to split node: %w", err)
		}
	}
	return nil
}

func allSameClassValue(instances []*WeightedInstance) bool {
	if len(instances) == 0 {
		return true
	}

	classValue := instances[0].Instance.ClassValue.Value().(string)
	for _, instance := range instances {
		if instance.Instance.ClassValue.Value().(string) != classValue {
			return false
		}
	}
	return true
}
