package tree

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"fmt"
)

func (t *Tree) Predict(instance *data.Instance) (string, error) {
	return t.RootNode.Predict(instance)
}

func (n *Node) Predict(instance *data.Instance) (string, error) {
	config.Logf("[Predict %d] Now at node %d\n", n.UniqId(), n.UniqId())
	if len(n.Children) == 0 {
		// find the majority class value
		config.Logf("[Predict %d] Reached leaf node, class: %s\n", n.UniqId(), n.LeafClass)
		return n.LeafClass, nil
	}

	attr := n.Children[0].Condition.Attr()
	val := instance.GetValueByAttr(attr)

	if !val.IsMissing() {
		for _, child := range n.Children {
			if child.Condition.IsMet(val) {
				config.Logf("[Predict %d] Value %v met condition <%s> to child node %d\n", n.UniqId(), val.Log(), child.Condition.Log(), child.UniqId())
				return child.Predict(instance)
			}
		}
		config.Logf("[Predict %d] Value %v mismatched all child nodes...\n", n.UniqId(), val.Log())
	}
	// If missing value, or no child is met, return the first prioritized child

	for _, child := range n.Children {
		if child.IsPrioritized {
			config.Logf("[Predict %d] Value %v goes along prioritized branch node %d...\n", n.UniqId(), val.Value(), child.UniqId())
			return child.Predict(instance)
		}
	}
	return "", fmt.Errorf("unknown error, cannot predict instance")
}
