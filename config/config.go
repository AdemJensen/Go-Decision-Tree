package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	ConsiderInvalidDataAsMissing bool `json:"consider_invalid_data_as_missing"`
}

func ReadConfig(filepath string) (*Config, error) {
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
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &config, nil
}
