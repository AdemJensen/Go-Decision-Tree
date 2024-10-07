package tree

import (
	"DecisionTree/data"
	"strings"
)

type Tree struct {
	Attributes []data.Attribute
	RootNode   *Node
}

func (t *Tree) Copy() *Tree {
	return &Tree{
		Attributes: t.Attributes,
		RootNode:   t.RootNode.Copy(),
	}
}

func (t *Tree) GetNodeCount() int {
	return t.RootNode.GetNodeCount()
}

func (t *Tree) GetLeafNodes() []*Node {
	return t.RootNode.GetLeafNodes()
}

func (t *Tree) GetMaxDepth() int {
	return t.RootNode.getMaxDepth(1)
}

type WeightedInstance struct {
	Instance *data.Instance
	Weight   float64 // 0~1, for missing value issues
}

func SumInstanceWeights(instances []*WeightedInstance) float64 {
	sum := 0.0
	for _, instance := range instances {
		sum += instance.Weight
	}
	return sum
}

func (w *WeightedInstance) CopyWithScale(scale float64) *WeightedInstance {
	return &WeightedInstance{
		Instance: w.Instance,
		Weight:   w.Weight * scale,
	}
}

type Node struct {
	Condition Condition
	Children  []*Node
	instances []*WeightedInstance // This is only valid during training

	IsPrioritized bool // When facing missing value, prioritize this node
	LeafClass     string

	uniqId int
}

var globalUniqId = 0

func (n *Node) UniqId() int {
	if n.uniqId == 0 {
		globalUniqId++
		n.uniqId = globalUniqId
	}
	return n.uniqId
}

func (n *Node) LogChildConditions() string {
	var conditions []string
	for _, child := range n.Children {
		conditions = append(conditions, "<"+child.Condition.Log()+">")
	}

	return "{" + strings.Join(conditions, ", ") + "}"
}

func (n *Node) Copy() *Node {
	var children []*Node
	for _, child := range n.Children {
		children = append(children, child.Copy())
	}
	return &Node{
		Condition:     n.Condition,
		Children:      children,
		instances:     n.instances,
		IsPrioritized: n.IsPrioritized,
		LeafClass:     n.LeafClass,
		uniqId:        n.uniqId,
	}
}

func (n *Node) GetNodeCount() int {
	count := 1
	for _, child := range n.Children {
		count += child.GetNodeCount()
	}
	return count
}

func (n *Node) GetLeafNodes() []*Node {
	if len(n.Children) == 0 {
		return []*Node{n}
	}
	var nodes []*Node
	for _, child := range n.Children {
		nodes = append(nodes, child.GetLeafNodes()...)
	}
	return nodes
}

func (n *Node) getMaxDepth(myDepth int) int {
	if len(n.Children) == 0 {
		return myDepth
	}
	var maxDepth int
	for _, child := range n.Children {
		depth := child.getMaxDepth(myDepth + 1)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}
