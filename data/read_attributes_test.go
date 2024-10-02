package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadAttributes(t *testing.T) {
	res, err := ReadAttributes("test_dataset/test_dataset.names")
	if err != nil {
		t.Errorf("Error reading attributes: %s", err)
		return
	}

	assert.NotNil(t, res.Class)
	assert.Equal(t, "Class", res.Class.Name(), "Attribute name not as expected")
	assert.Equal(t, []string{"OK", "Not OK"}, res.Class.AcceptedValues, "Class.AcceptedValues")
	assert.Equal(t, 4, len(res.Attributes), "Number of attributes")
	assert.Equal(t, "Attribute1", res.Attributes[0].Name(), "Attribute1 name")
	assert.Equal(t, "Attribute2", res.Attributes[1].Name(), "Attribute2 name")
	assert.Equal(t, "Attribute 3", res.Attributes[2].Name(), "Attribute3 name")
	assert.Equal(t, "Attribute 4", res.Attributes[3].Name(), "Attribute4 name")
	assert.Equal(t, Continuous, res.Attributes[0].Type(), "Attribute1 type")
	assert.Equal(t, Continuous, res.Attributes[1].Type(), "Attribute2 type")
	assert.Equal(t, Nominal, res.Attributes[2].Type(), "Attribute3 type")
	assert.Equal(t, []string{"A", "B", "C"}, res.Attributes[2].(*NominalAttribute).AcceptedValues, "Attribute3 list")
	assert.Equal(t, Nominal, res.Attributes[3].Type(), "Attribute4 type")
	assert.Equal(t, []string{"D", "E", "F"}, res.Attributes[3].(*NominalAttribute).AcceptedValues, "Attribute4 list")
}
