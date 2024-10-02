package data

import (
	"DecisionTree/config"
	"fmt"
)

type AttributeType string

const (
	Unknown    AttributeType = ""
	Continuous AttributeType = "continuous"
	Nominal    AttributeType = "nominal"
)

type Attribute interface {
	Name() string
	Type() AttributeType
	Parse(conf *config.Config, value string) (Value, error)
}

type ContinuousAttribute struct {
	name string
}

func (c *ContinuousAttribute) Name() string {
	return c.name
}

func (c *ContinuousAttribute) Type() AttributeType {
	return Continuous
}

func (c *ContinuousAttribute) Parse(conf *config.Config, value string) (Value, error) {
	res, err := newContinuousValue(conf, c, value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse continuous value: %w", err)
	}
	return res, nil
}

type NominalAttribute struct {
	name           string
	AcceptedValues []string
}

func (n *NominalAttribute) Name() string {
	return n.name
}

func (n *NominalAttribute) Type() AttributeType {
	return Nominal
}

func (n *NominalAttribute) Parse(conf *config.Config, value string) (Value, error) {
	res, err := newNominalValue(conf, n, value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nominal value: %w", err)
	}
	return res, nil
}
