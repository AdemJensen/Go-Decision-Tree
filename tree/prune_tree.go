package tree

import (
	"DecisionTree/config"
	"DecisionTree/data"
	"fmt"

	"github.com/gosuri/uiprogress"
)

func postPruneTree(conf *config.Config, tree *Tree, trainData *data.ValueTable) error {
	// Get all prune-ready nodes
	pruneReadyNodes := getPruneReadyNodes(tree.RootNode)
	bar := uiprogress.AddBar(len(pruneReadyNodes)).PrependFunc(func(b *uiprogress.Bar) string {
		return "Post-pruning"
	}).AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d/%d", b.Current(), b.Total)
	})
	// build reverse mapping
	reverseMapping := buildReverseMapping(tree.RootNode)
	// get instances related to each node's prediction
	instancesMapping := getInstancesRelatedToNodePrediction(tree.RootNode, trainData.Instances)

	// for each leaf node
	for len(pruneReadyNodes) > 0 {
		targetNode := pruneReadyNodes[0]
		pruneReadyNodes = pruneReadyNodes[1:]

		// get err count related to this node
		oldInfo, err := testRunNode(targetNode, instancesMapping[targetNode.UniqId()])
		if err != nil {
			return fmt.Errorf("failed to test run node (prior): %w", err)
		}

		// try how much error will be reduced if we prune this node
		savedChildren := targetNode.Children
		targetNode.Children = nil
		err = postProcessNode(conf, targetNode)
		if err != nil {
			return fmt.Errorf("failed to post process node: %w", err)
		}

		// calculate its new pessimistic error
		newInfo, err := testRunNode(targetNode, instancesMapping[targetNode.UniqId()])
		if err != nil {
			return fmt.Errorf("failed to test run node (post): %w", err)
		}
		if -(newInfo.PessimisticError - oldInfo.PessimisticError) < conf.MinPostPruneGeneralizationErrorDecrease {
			// if the error is not decreased, revert the prune
			targetNode.Children = savedChildren
		} else {
			// if the error is decreased, add its parent to the prune ready nodes
			if isNodePruneReady(reverseMapping[targetNode.UniqId()]) {
				pruneReadyNodes = append(pruneReadyNodes, reverseMapping[targetNode.UniqId()])
				bar.Total++
			}
			config.Logf("[Post-Prune] Pruned node %d, Pessimistic Error: %.6f%% -> %.6f%% (%.6f%%), Error Nodes: %d -> %d (%+d), Leaf Nodes: %d -> %d (%+d)",
				targetNode.UniqId(), oldInfo.PessimisticError*100, newInfo.PessimisticError*100, (newInfo.PessimisticError-oldInfo.PessimisticError)*100,
				oldInfo.ErrorCount, newInfo.ErrorCount, newInfo.ErrorCount-oldInfo.ErrorCount,
				len(savedChildren), 1, 1-len(savedChildren))
		}
		bar.Incr()
	}

	return nil
}

func getLeafNodes(n *Node) []*Node {
	if len(n.Children) == 0 {
		return []*Node{n}
	}
	var nodes []*Node
	for _, child := range n.Children {
		nodes = append(nodes, getLeafNodes(child)...)
	}
	return nodes
}

// getPruneReadyNodes returns all nodes that are ready to be pruned
// A node is ready to be pruned if all its children are leaf nodes
func getPruneReadyNodes(n *Node) []*Node {
	if len(n.Children) == 0 {
		return nil
	}
	if isNodePruneReady(n) {
		return []*Node{n}
	} else {
		var nodes []*Node
		for _, child := range n.Children {
			nodes = append(nodes, getPruneReadyNodes(child)...)
		}
		return nodes
	}
}

func isNodePruneReady(node *Node) bool {
	if node == nil || len(node.Children) == 0 {
		return false
	}
	for _, child := range node.Children {
		if len(child.Children) != 0 {
			return false
		}
	}
	return true
}

func buildReverseMapping(node *Node) map[int]*Node {
	mapping := make(map[int]*Node)
	for _, child := range node.Children {
		mapping[child.UniqId()] = node
		for k, v := range buildReverseMapping(child) {
			mapping[k] = v
		}
	}
	return mapping
}

func getInstancesRelatedToNodePrediction(node *Node, instances []*data.Instance) map[int][]*data.Instance {
	mapping := make(map[int][]*data.Instance)
	nextLayerMapping := make(map[int][]*data.Instance)
	for _, instance := range instances {
		mapping[node.UniqId()] = append(mapping[node.UniqId()], instance)
		relatedChild := node.GetRelatedChild(instance)
		if relatedChild != nil {
			nextLayerMapping[relatedChild.UniqId()] = append(nextLayerMapping[relatedChild.UniqId()], instance)
		}
	}

	for _, child := range node.Children {
		for k, v := range getInstancesRelatedToNodePrediction(child, nextLayerMapping[child.UniqId()]) {
			mapping[k] = append(mapping[k], v...)
		}
	}
	return mapping
}
