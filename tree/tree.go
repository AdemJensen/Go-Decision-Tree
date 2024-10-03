package tree

import "DecisionTree/data"

type Tree struct {
	RootNode *Node
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
}
