package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyAppName       = errors.New("app name cannot be empty")
	ErrInvalidPort        = errors.New("port must be between 1 and 65535")
	ErrInvalidEnvironment = errors.New("environment must be development, testing or production")
)

type Config struct {
	AppName     string          `json:"app_name"`
	Environment string          `json:"environment"`
	Port        int             `json:"port"`
	Features    map[string]bool `json:"features,omitempty"`
	Hosts       []string        `json:"hosts,omitempty"`
}

// Validate 检查业务规则；JSON 能解码不代表配置可用。
func (config Config) Validate() error {
	if strings.TrimSpace(config.AppName) == "" {
		return ErrEmptyAppName
	}
	if config.Port < 1 || config.Port > 65535 {
		return ErrInvalidPort
	}
	switch config.Environment {
	case "development", "testing", "production":
		return nil
	default:
		return ErrInvalidEnvironment
	}
}

func DecodeConfig(data []byte) (Config, error) {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func EncodeConfig(config Config) ([]byte, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode config: %w", err)
	}
	return data, nil
}

func main() {
	input := []byte(`{
		"app_name": "gateway",
		"environment": "development",
		"port": 8080,
		"features": {"metrics": true, "debug": false},
		"hosts": ["localhost", "127.0.0.1"]
	}`)

	config, err := DecodeConfig(input)
	if err != nil {
		fmt.Println("配置无效:", err)
		return
	}

	output, err := EncodeConfig(config)
	if err != nil {
		fmt.Println("编码失败:", err)
		return
	}
	fmt.Println(string(output))
}
