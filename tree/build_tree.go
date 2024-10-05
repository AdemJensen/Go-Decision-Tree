package tree

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"fmt"
)

func BuildTree(conf *config.Config, valueTable *data.ValueTable) (*Tree, error) {
	// wash data without class values
	var instances []*WeightedInstance
	for _, instance := range valueTable.Instances {
		if instance.ClassValue.IsMissing() {
			continue
		}
		instances = append(instances, &WeightedInstance{
			Instance: instance,
			Weight:   1,
		})
	}

	tree := &Tree{
		RootNode: &Node{
			instances: instances,
		},
	}

	// split node
	err := splitNode(conf, 1, tree.RootNode)
	if err != nil {
		return nil, fmt.Errorf("failed to split node: %w", err)
	}

	// post prune tree
	err = postPruneTree(conf, tree)
	if err != nil {
		return nil, fmt.Errorf("failed to post prune tree: %w", err)
	}

	// post process tree
	err = postProcessTree(conf, tree)
	if err != nil {
		return nil, fmt.Errorf("failed to post process tree: %w", err)
	}

	return tree, nil
}
