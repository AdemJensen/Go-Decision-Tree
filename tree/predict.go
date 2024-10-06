package tree

import (
	"DecisionTree/data"
	"fmt"
)

func (t *Tree) Predict(instance *data.Instance) (string, error) {
	return t.RootNode.Predict(instance)
}

func (n *Node) Predict(instance *data.Instance) (string, error) {
	if len(n.Children) == 0 {
		// find the majority class value
		return n.LeafClass, nil
	}

	attr := n.Children[0].Condition.Attr()
	val := instance.GetValueByAttr(attr)

	if !val.IsMissing() {
		for _, child := range n.Children {
			if child.Condition.IsMet(val) {
				return child.Predict(instance)
			}
		}
	}
	// If missing value, or no child is met, return the first prioritized child

	for _, child := range n.Children {
		if child.IsPrioritized {
			return child.Predict(instance)
		}
	}
	return "", fmt.Errorf("unknown error, cannot predict instance")
}
