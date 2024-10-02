package data

import (
	"DecisionTree/config"
	"testing"
)

func TestReadValues(t *testing.T) {
	attrTable, err := ReadAttributes("test_dataset/test_dataset.names")
	if err != nil {
		t.Errorf("Error reading attributes: %s", err)
		return
	}
	res, err := ReadValues(&config.Config{}, attrTable, "test_dataset/test_dataset.data")
	if err != nil {
		t.Errorf("Error reading data: %s", err)
		return
	}

	t.Log(res.String())
}
