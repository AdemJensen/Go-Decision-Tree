package data

import (
	"DecisionTree/config"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type ValueTable struct {
	Instances []*Instance
}

func (v *ValueTable) String() string {
	var sb strings.Builder
	for _, instance := range v.Instances {
		sb.WriteString(fmt.Sprintf("%s\n", instance.String()))
	}
	return sb.String()
}

type Instance struct {
	AttributeValues []Value
	ClassValue      *NominalValue
}

func (i *Instance) GetValueByAttr(attr Attribute) Value {
	for _, value := range i.AttributeValues {
		if value.Attribute().Name() == attr.Name() {
			return value
		}
	}
	return nil
}

func (i *Instance) String() string {
	var sb strings.Builder
	for _, value := range i.AttributeValues {
		if value.IsMissing() {
			sb.WriteString("?, ")
		} else {
			sb.WriteString(fmt.Sprintf("%v, ", value.Value()))
		}
	}
	if i.ClassValue.IsMissing() {
		sb.WriteString("?")
	} else {
		sb.WriteString(fmt.Sprintf("%v", i.ClassValue.Value()))
	}
	return sb.String()
}

func ReadValues(conf *config.Config, attrTable *AttributeTable, filepath string) (*ValueTable, error) {
	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	table := &ValueTable{}

	// Read file line by line
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		instance, err := handleInstanceLine(conf, attrTable, line)
		if err != nil {
			log.Printf("Error parsing line %d: %s\n", lineNo, err)
			continue
		}
		if instance == nil {
			continue
		}

		table.Instances = append(table.Instances, instance)
	}
	return table, nil
}

func handleInstanceLine(conf *config.Config, attrTable *AttributeTable, line string) (*Instance, error) {
	line = strings.TrimSpace(line)      // Remove leading and trailing spaces
	line = strings.TrimRight(line, ".") // Remove trailing period

	// Ignore empty lines
	if len(line) == 0 {
		return nil, nil
	}
	// Ignore comment lines
	if len(line) == 0 || strings.HasPrefix(line, "|") {
		return nil, nil
	}

	dataValues := strings.Split(line, ",")

	// If data is not sufficient, return error
	if len(dataValues) < len(attrTable.Attributes) {
		return nil, fmt.Errorf("insufficient data values, expected %d, got %d", len(attrTable.Attributes), len(dataValues))
	}

	// Parse each value
	instance := &Instance{}
	for i, attr := range attrTable.Attributes {
		dataValue := dataValues[i]
		value, err := attr.Parse(conf, dataValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value '%s': %w", dataValue, err)
		}
		instance.AttributeValues = append(instance.AttributeValues, value)
	}

	// Parse class value
	if len(dataValues) >= len(attrTable.Attributes)+1 {
		classValue, err := attrTable.Class.Parse(conf, dataValues[len(dataValues)-1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse class value '%s': %w", dataValues[len(dataValues)-1], err)
		}
		instance.ClassValue = classValue.(*NominalValue)
	} else {
		newVal, err := newNominalValue(conf, attrTable.Class, "?") // create a default class value, avoid nil pointer
		if err != nil {
			return nil, fmt.Errorf("failed to create default missing class value: %w", err)
		}
		instance.ClassValue = newVal.(*NominalValue)
	}

	return instance, nil
}
