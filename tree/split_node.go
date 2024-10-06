package tree

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"fmt"
)

func splitNode(conf *config.Config, level int, node *Node) error {
	node.UniqId() // ensure node has a uniq id

	// if reach max depth, stop split
	if level >= conf.MaxDepth {
		config.Logf("[Train %d] [Level %d] Reach max depth, stop split", node.UniqId(), level)
		return nil
	}

	// if reach min samples split, stop split
	if len(node.instances) < conf.MinSamplesSplit {
		config.Logf("[Train %d] [Level %d] Reach min samples split (n_instance=%d, min_samples_split=%d), stop split", node.UniqId(), level, len(node.instances), conf.MinSamplesSplit)
		return nil
	}

	// if all instances have the same class value, stop split
	if allSameClassValue(node.instances) {
		config.Logf("[Train %d] [Level %d] All instances have the same class value, stop split", node.UniqId(), level)
		return nil
	}

	if node.UniqId() == 13 {
		println()
	}

	// for each attribute, try to split node, find the best split
	// Initial best split: do not split. Impurity is the impurity of the node
	var (
		bestSplitChildren []*Node // empty node list, means do not split
		bestSplitGain     = 0.0
		nodeEntropy       = calculateEntropy(node.instances, 0, nil)
		//bestSplitIndex    int
	)
	for i, value := range node.instances[0].Instance.AttributeValues {
		if node.UniqId() == 13 && i == 1 {
			println()
		}
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
				//bestSplitIndex = i
			}
		case data.Nominal:
			bestNominalSplit, bestNominalGain, err := splitInstancesByNominalAttr(conf, nodeEntropy, i, node.instances)
			if err != nil {
				return fmt.Errorf("failed to split instances by nominal attribute: %w", err)
			}
			if len(bestNominalSplit) > 0 && bestNominalGain > bestSplitGain {
				bestSplitGain = bestNominalGain
				bestSplitChildren = buildNodeListFromNominalSplit(attribute, bestNominalSplit)
				//bestSplitIndex = i
			}
		}
	}

	// if no split, stop split
	if len(bestSplitChildren) == 0 || bestSplitGain == 0 {
		config.Logf("[Train %d] [Level %d] No split found, stop split", node.UniqId(), level)
		return nil
	}

	if bestSplitGain < conf.MinImpurityDecrease {
		config.Logf("[Train %d] [Level %d] Best split gain did not meet threshold (%.2f vs %.2f), stop split", node.UniqId(), level, bestSplitGain, conf.MinImpurityDecrease)
		return nil
	}

	// split node
	node.Children = bestSplitChildren

	// DEBUG
	//if len(node.Children) == 1 {
	//	println()
	//	nodeEntropy = calculateEntropy(node.instances, 0, nil)
	//	_, bestSplitGain, _ = splitInstancesByNominalAttr(conf, nodeEntropy, bestSplitIndex, node.instances)
	//}

	config.Logf("[Train %d] [Level %d] Split node by condition %s, gain=%f, n_instance=%d", node.UniqId(), level, node.LogChildConditions(), bestSplitGain, len(node.instances))
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
