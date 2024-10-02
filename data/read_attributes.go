package data

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type AttributeTable struct {
	Attributes []Attribute
	Class      *NominalAttribute
}

// ReadAttributes reads the attributes from the file and returns them.
// The file follows the format:
// <attribute name>: continuous.
// or:
// <attribute name>: <V1>, <V2>, <V3>.
// The first attribute is the class attribute.
// Empty lines or lines starting with '|' are ignored.
func ReadAttributes(filepath string) (*AttributeTable, error) {
	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	table := &AttributeTable{}

	// Read file line by line
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		attr, err := handleAttributeLine(line)
		if err != nil {
			log.Printf("Error parsing line %d: %s\n", lineNo, err)
			continue
		}
		if attr == nil {
			continue
		}

		if table.Class == nil {
			nominalAttr, ok := attr.(*NominalAttribute)
			if !ok {
				return nil, fmt.Errorf("first attribute is considered class attribute, it must be nominal")
			}
			if nominalAttr.Name() == "" {
				nominalAttr.name = "Class"
			}
			table.Class = nominalAttr
		} else {
			if attr.Name() == "" {
				return nil, fmt.Errorf("non-class attribute on line %d does not have a name", lineNo)
			}
			table.Attributes = append(table.Attributes, attr)
		}
	}

	return table, nil
}

func handleAttributeLine(line string) (Attribute, error) {
	line = strings.TrimSpace(line)      // Remove leading and trailing spaces
	line = strings.TrimRight(line, ".") // Remove trailing period

	// Ignore comment lines
	if len(line) == 0 || strings.HasPrefix(line, "|") {
		return nil, nil
	}

	// Check and handle line data
	attrName := "" // name is optional
	attrContent := line
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 2) // 分割行内容为 [name] 和 剩余部分
		attrName = strings.TrimSpace(parts[0])
		attrContent = strings.TrimSpace(parts[1])
	}

	// Check if the attribute is nominal or continuous
	switch {
	case attrContent == "continuous":
		return &ContinuousAttribute{name: attrName}, nil
	default:
		// Parse nominal attribute values
		values := strings.Split(attrContent, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		return &NominalAttribute{name: attrName, AcceptedValues: values}, nil
	}
}
