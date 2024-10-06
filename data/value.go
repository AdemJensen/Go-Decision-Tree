package data

import (
	"DecisionTree/config"
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	Attribute() Attribute
	IsMissing() bool
	Value() interface{}
	Log() string
}

type ContinuousValue struct {
	attr      Attribute
	isMissing bool
	value     float64
}

func newContinuousValue(conf *config.Config, attr *ContinuousAttribute, value string) (Value, error) {
	value = strings.TrimSpace(value)
	if value == "?" {
		return &ContinuousValue{
			attr:      attr,
			isMissing: true,
		}, nil
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		if conf.ConsiderInvalidDataAsMissing {
			return &ContinuousValue{
				attr:      attr,
				isMissing: true,
			}, nil
		}
		return nil, fmt.Errorf("failed to parse value '%s' to float: %w", value, err)
	}
	return &ContinuousValue{
		attr:  attr,
		value: val,
	}, nil
}

func (c ContinuousValue) Attribute() Attribute {
	return c.attr
}

func (c ContinuousValue) IsMissing() bool {
	return c.isMissing
}

func (c ContinuousValue) Value() interface{} {
	return c.value
}

func (c ContinuousValue) Log() string {
	if c.isMissing {
		return "<missing>"
	}
	return fmt.Sprintf("%f", c.value)
}

func newNominalValue(conf *config.Config, attr *NominalAttribute, value string) (Value, error) {
	value = strings.TrimSpace(value)
	if value == "?" {
		return &NominalValue{
			attr:      attr,
			isMissing: true,
		}, nil
	}

	// check if value is in accepted values
	for _, acceptedValue := range attr.AcceptedValues {
		if value == acceptedValue {
			return &NominalValue{
				attr:  attr,
				value: value,
			}, nil
		}
	}
	// value is not in accepted values
	if conf.ConsiderInvalidDataAsMissing {
		return &NominalValue{
			attr:      attr,
			isMissing: true,
		}, nil
	}
	return nil, fmt.Errorf("value '%s' is not in accepted values", value)
}

type NominalValue struct {
	attr      Attribute
	isMissing bool
	value     string
}

func (n NominalValue) Attribute() Attribute {
	return n.attr
}

func (n NominalValue) IsMissing() bool {
	return n.isMissing
}

func (n NominalValue) Value() interface{} {
	return n.value
}

func (n NominalValue) Log() string {
	if n.isMissing {
		return "<missing>"
	}
	return n.value
}
