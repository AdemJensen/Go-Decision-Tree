package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	ConsiderInvalidDataAsMissing bool    `json:"consider_invalid_data_as_missing"`
	MaxDepth                     int     `json:"max_depth"`
	MinSamplesSplit              int     `json:"min_samples_split"`
	MinSamplesLeaf               int     `json:"min_samples_leaf"`
	MinImpurityDecrease          float64 `json:"min_impurity_decrease"`

	// For nominal attribute, if the number of accepted values is less than this value, use brute-force to find
	// the best split. If not, we will first join the values with fewer instances until the number of values is
	// less than or equal to this value, then perform brute-force.
	// This value must >= 2.
	MaxNominalBruteForceScale int `json:"max_nominal_brute_force_scale"`

	VerboseLog bool   `json:"verbose_log"`
	LogFile    string `json:"log_file"`
}

func Logf(format string, v ...interface{}) {
	if Conf.VerboseLog {
		log.Printf(format, v...)
	}
}

var Conf *Config

func init() {
	confPath := os.Getenv("CONF_PATH")
	if confPath == "" {
		confPath = "config.json"
	}
	conf, err := ReadConfig(confPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}
	Conf = conf

	// log file
	if Conf.LogFile != "" {
		file, err := os.OpenFile(Conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			panic(fmt.Errorf("failed to open log file: %w", err))
		}
		log.SetOutput(file)
	}
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
