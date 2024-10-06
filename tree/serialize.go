package tree

import (
	"DecisionTree/data"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadTreeFromFile(filepath string) (*Tree, error) {
	// read json content from file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 反序列化 JSON 内容
	var pt PersistentTree
	if err := json.Unmarshal(bytes, &pt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return pt.ToTree(), nil
}

func WriteTreeToFile(tree *Tree, filepath string) error {
	pt := NewPersistentTree(tree)
	bytes, err := json.Marshal(pt)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if _, err := file.Write(bytes); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

type PersistentTree struct {
	Attributes []*data.PersistentAttribute `json:"attributes"`
	RootNode   *PersistentNode             `json:"root_node"`
}

func NewPersistentTree(tree *Tree) *PersistentTree {
	var attrList []*data.PersistentAttribute
	for _, attr := range tree.Attributes {
		attrList = append(attrList, data.NewPersistentAttribute(attr))
	}
	return &PersistentTree{
		Attributes: attrList,
		RootNode:   NewPersistentNode(attrList, tree.RootNode),
	}
}

func (p *PersistentTree) ToTree() *Tree {
	if p == nil {
		return nil
	}
	var attrList []data.Attribute
	for _, attr := range p.Attributes {
		attrInst, _ := attr.ToAttribute()
		attrList = append(attrList, attrInst)
	}
	return &Tree{
		Attributes: attrList,
		RootNode:   p.RootNode.ToNode(attrList),
	}
}

type PersistentNode struct {
	UniqId        int `json:"uniq_id"`
	Condition     *PersistentCondition
	Children      []*PersistentNode
	IsPrioritized bool   `json:"is_prioritized,omitempty"`
	LeafClass     string `json:"leaf_class,omitempty"`
}

func NewPersistentNode(attrList []*data.PersistentAttribute, node *Node) *PersistentNode {
	pNode := &PersistentNode{
		UniqId:        node.UniqId(),
		Condition:     NewPersistentCondition(attrList, node.Condition),
		Children:      nil,
		IsPrioritized: node.IsPrioritized,
		LeafClass:     node.LeafClass,
	}
	for _, child := range node.Children {
		pNode.Children = append(pNode.Children, NewPersistentNode(attrList, child))
	}
	return pNode
}

func (p *PersistentNode) ToNode(attrList []data.Attribute) *Node {
	if p == nil {
		return nil
	}
	node := &Node{
		Condition:     p.Condition.ToCondition(attrList),
		Children:      nil,
		instances:     nil,
		IsPrioritized: p.IsPrioritized,
		LeafClass:     p.LeafClass,
		uniqId:        p.UniqId,
	}
	for _, child := range p.Children {
		node.Children = append(node.Children, child.ToNode(attrList))
	}
	return node
}

type PersistentCondition struct {
	ConditionType  ConditionType `json:"condition_type"`
	AttrId         int           `json:"attr_id"`
	UpperValue     float64       `json:"upper_value,omitempty"`
	LowerValue     float64       `json:"lower_value,omitempty"`
	AcceptedValues []string      `json:"accepted_values,omitempty"`
}

func getIdFromAttrList(attrList []*data.PersistentAttribute, attribute data.Attribute) int {
	for i, attr := range attrList {
		if attr.Name == attribute.Name() {
			return i
		}
	}
	return -1
}

func NewPersistentCondition(attrList []*data.PersistentAttribute, condition Condition) *PersistentCondition {
	switch c := condition.(type) {
	case *ContinuousCondition:
		return &PersistentCondition{
			ConditionType: c.conditionType,
			AttrId:        getIdFromAttrList(attrList, c.Attr()),
			UpperValue:    c.upperValue,
			LowerValue:    c.lowerValue,
		}
	case *NominalCondition:
		return &PersistentCondition{
			ConditionType:  c.conditionType,
			AttrId:         getIdFromAttrList(attrList, c.Attr()),
			AcceptedValues: c.acceptedValues,
		}
	default:
		return nil
	}
}

func (p *PersistentCondition) ToCondition(attrList []data.Attribute) Condition {
	if p == nil {
		return nil
	}
	switch p.ConditionType {
	case LessThan:
		return newLessThanCondition(attrList[p.AttrId], p.UpperValue)
	case GreaterThanEq:
		return newGreaterThanEqCondition(attrList[p.AttrId], p.LowerValue)
	// Range is not supported
	case IsOneOf:
		return newIsOneOfCondition(attrList[p.AttrId], p.AcceptedValues)
	default:
		return nil
	}
}
