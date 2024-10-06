package data

import "fmt"

type PersistentAttribute struct {
	Name string        `json:"name"`
	Type AttributeType `json:"type"`
}

func NewPersistentAttribute(attr Attribute) *PersistentAttribute {
	return &PersistentAttribute{
		Name: attr.Name(),
		Type: attr.Type(),
	}
}

func (p *PersistentAttribute) ToAttribute() (Attribute, error) {
	switch p.Type {
	case Continuous:
		return &ContinuousAttribute{name: p.Name}, nil
	case Nominal:
		return &NominalAttribute{name: p.Name}, nil
	default:
		return nil, fmt.Errorf("unknown attribute type: %s", p.Type)
	}
}
