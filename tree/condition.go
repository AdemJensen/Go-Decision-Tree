package tree

import "DecisionTree/data"

type ConditionType string

const (
	LessThan      ConditionType = "<"         // continuous
	GreaterThanEq ConditionType = ">="        // continuous
	Range         ConditionType = "range"     // continuous, min < value <= max
	IsOneOf       ConditionType = "is_one_of" // Nominal
)

type Condition interface {
	Type() ConditionType
	Attr() data.Attribute
	IsMet(value data.Value) bool
}

type ContinuousCondition struct {
	conditionType ConditionType
	attr          data.Attribute
	upperValue    float64
	lowerValue    float64
}

func newLessThanCondition(attr data.Attribute, upperValue float64) *ContinuousCondition {
	return &ContinuousCondition{
		conditionType: LessThan,
		attr:          attr,
		upperValue:    upperValue,
	}
}

func newGreaterThanEqCondition(attr data.Attribute, lowerValue float64) *ContinuousCondition {
	return &ContinuousCondition{
		conditionType: GreaterThanEq,
		attr:          attr,
		lowerValue:    lowerValue,
	}
}

func (c *ContinuousCondition) Type() ConditionType {
	return c.conditionType
}

func (c *ContinuousCondition) Attr() data.Attribute {
	return c.attr
}

func (c *ContinuousCondition) IsMet(value data.Value) bool {
	v := value.Value().(float64)
	switch c.conditionType {
	case LessThan:
		return v < c.upperValue
	case GreaterThanEq:
		return v >= c.lowerValue
	case Range:
		return c.lowerValue < v && v <= c.upperValue
	default:
		return false
	}
}

type NominalCondition struct {
	conditionType  ConditionType
	attr           data.Attribute
	acceptedValues []string
}

func newIsOneOfCondition(attr data.Attribute, acceptedValues []string) *NominalCondition {
	return &NominalCondition{
		conditionType:  IsOneOf,
		attr:           attr,
		acceptedValues: acceptedValues,
	}
}

func (n *NominalCondition) Type() ConditionType {
	return n.conditionType
}

func (n *NominalCondition) Attr() data.Attribute {
	return n.attr
}

func (n *NominalCondition) IsMet(value data.Value) bool {
	v := value.Value().(string)
	for _, acceptedValue := range n.acceptedValues {
		if v == acceptedValue {
			return true
		}
	}
	return false
}
