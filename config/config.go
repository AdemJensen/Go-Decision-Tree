package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	ConsiderInvalidDataAsMissing bool `json:"consider_invalid_data_as_missing"`
	MaxDepth                     int  `json:"max_depth"`
	MinSamplesSplit              int  `json:"min_samples_split"`
	MinSamplesLeaf               int  `json:"min_samples_leaf"`

	// For nominal attribute, if the number of accepted values is less than this value, use brute-force to find
	// the best split. If not, we will first join the values with fewer instances until the number of values is
	// less than or equal to this value, then perform brute-force.
	// This value must >= 2.
	MaxNominalBruteForceScale int `json:"max_nominal_brute_force_scale"`
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
