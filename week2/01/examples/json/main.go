package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyAppName = errors.New("app name cannot be empty")
	ErrInvalidPort  = errors.New("port must be between 1 and 65535")
)

type Config struct {
	AppName  string          `json:"app_name"`
	Port     int             `json:"port"`
	Features map[string]bool `json:"features,omitempty"`
	Hosts    []string        `json:"hosts,omitempty"`
}

func DecodeConfig(data []byte) (Config, error) {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	if strings.TrimSpace(config.AppName) == "" {
		return Config{}, ErrEmptyAppName
	}
	if config.Port < 1 || config.Port > 65535 {
		return Config{}, ErrInvalidPort
	}
	return config, nil
}

func main() {
	input := []byte(`{
		"app_name": "gateway",
		"port": 8080,
		"features": {"metrics": true, "debug": false},
		"hosts": ["localhost", "127.0.0.1"]
	}`)

	config, err := DecodeConfig(input)
	if err != nil {
		fmt.Println("配置无效:", err)
		return
	}

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("编码失败:", err)
		return
	}
	fmt.Println(string(output))
}
