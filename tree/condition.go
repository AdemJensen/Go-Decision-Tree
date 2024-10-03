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
	IsMet(value data.Value) bool
}

type ContinuousCondition struct {
	conditionType ConditionType
	upperValue    float64
	lowerValue    float64
}

func newLessThanCondition(upperValue float64) *ContinuousCondition {
	return &ContinuousCondition{
		conditionType: LessThan,
		upperValue:    upperValue,
	}
}

func newGreaterThanEqCondition(lowerValue float64) *ContinuousCondition {
	return &ContinuousCondition{
		conditionType: GreaterThanEq,
		lowerValue:    lowerValue,
	}
}

func newRangeCondition(lowerValue, upperValue float64) *ContinuousCondition {
	return &ContinuousCondition{
		conditionType: Range,
		upperValue:    upperValue,
		lowerValue:    lowerValue,
	}
}

func (c *ContinuousCondition) Type() ConditionType {
	return c.conditionType
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
	acceptedValues []string
}

func newIsOneOfCondition(acceptedValues []string) *NominalCondition {
	return &NominalCondition{
		conditionType:  IsOneOf,
		acceptedValues: acceptedValues,
	}
}

func (n *NominalCondition) Type() ConditionType {
	return n.conditionType
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
