package tree

import (
	"DecisionTree/config"
)

func postProcessTree(conf *config.Config, tree *Tree) error {
	// for each node, calculate its majority class
	return postProcessNode(conf, tree.RootNode)
}

func postProcessNode(_ *config.Config, node *Node) error {
	// if is leaf node, calculate its majority class
	if len(node.Children) == 0 {
		classFrequency := make(map[string]float64)
		for _, ins := range node.instances {
			classValue := ins.Instance.ClassValue.Value().(string)
			classFrequency[classValue] += ins.Weight
		}
		var (
			maxFrequency      float64
			maxFrequencyClass string
		)
		for c, f := range classFrequency {
			if f > maxFrequency {
				maxFrequency = f
				maxFrequencyClass = c
			}
		}
		node.LeafClass = maxFrequencyClass
	} else {
		for _, child := range node.Children {
			if err := postProcessNode(nil, child); err != nil {
				return err
			}
		}
	}

	return nil
}
