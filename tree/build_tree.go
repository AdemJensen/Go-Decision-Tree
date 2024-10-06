package tree

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"fmt"

	"github.com/gosuri/uiprogress"
)

func BuildTree(conf *config.Config, valueTable *data.ValueTable) (*Tree, error) {
	uiprogress.Start()
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

	if len(instances) == 0 {
		return nil, fmt.Errorf("no valid instances")
	}

	var attributes []data.Attribute
	for _, attr := range instances[0].Instance.AttributeValues {
		attributes = append(attributes, attr.Attribute())
	}

	tree := &Tree{
		Attributes: attributes,
		RootNode: &Node{
			instances: instances,
		},
	}

	// split node
	bar := uiprogress.AddBar(1).PrependFunc(func(b *uiprogress.Bar) string {
		return "Building Nodes"
	}).AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d/%d", b.Current(), b.Total)
	})
	err := splitNode(bar, conf, 1, tree.RootNode)
	if err != nil {
		return nil, fmt.Errorf("failed to split node: %w", err)
	}

	// post process tree
	err = postProcessTree(conf, tree)
	if err != nil {
		return nil, fmt.Errorf("failed to post process tree: %w", err)
	}

	// post prune tree
	err = postPruneTree(conf, tree, valueTable)
	if err != nil {
		return nil, fmt.Errorf("failed to post prune tree: %w", err)
	}

	return tree, nil
}
