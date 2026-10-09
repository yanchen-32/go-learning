package main

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestConfigRoundTrip(t *testing.T) {
	input := []byte(`{"app_name":"gateway","environment":"production","port":8080,"features":{"metrics":true,"debug":false},"hosts":["localhost","127.0.0.1"]}`)
	config, err := DecodeConfig(input)
	if err != nil {
		t.Fatal(err)
	}
	want := Config{AppName: "gateway", Environment: "production", Port: 8080, Features: map[string]bool{"metrics": true, "debug": false}, Hosts: []string{"localhost", "127.0.0.1"}}
	if !reflect.DeepEqual(config, want) {
		t.Fatalf("config = %+v, want %+v", config, want)
	}
	output, err := EncodeConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeConfig(output)
	if err != nil || !reflect.DeepEqual(decoded, want) {
		t.Fatalf("round trip = %+v, %v", decoded, err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(output, &fields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"app_name", "environment", "port", "features", "hosts"} {
		if _, ok := fields[key]; !ok {
			t.Errorf("JSON tag missing for %q", key)
		}
	}
}

func TestDecodeConfigInvalid(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		want  error
	}{
		{"syntax", `{"app_name":`, nil},
		{"wrong port type", `{"app_name":"app","environment":"testing","port":"8080"}`, nil},
		{"wrong features type", `{"app_name":"app","environment":"testing","port":8080,"features":[]}`, nil},
		{"wrong hosts type", `{"app_name":"app","environment":"testing","port":8080,"hosts":{}}`, nil},
		{"missing name", `{"environment":"testing","port":8080}`, ErrEmptyAppName},
		{"blank name", `{"app_name":"  ","environment":"testing","port":8080}`, ErrEmptyAppName},
		{"missing port", `{"app_name":"app","environment":"testing"}`, ErrInvalidPort},
		{"zero port", `{"app_name":"app","environment":"testing","port":0}`, ErrInvalidPort},
		{"negative port", `{"app_name":"app","environment":"testing","port":-1}`, ErrInvalidPort},
		{"large port", `{"app_name":"app","environment":"testing","port":65536}`, ErrInvalidPort},
		{"missing environment", `{"app_name":"app","port":8080}`, ErrInvalidEnvironment},
		{"invalid environment", `{"app_name":"app","environment":"other","port":8080}`, ErrInvalidEnvironment},
		{"null", `null`, ErrEmptyAppName},
	} {
		t.Run(tt.name, func(t *testing.T) {
			config, err := DecodeConfig([]byte(tt.input))
			if err == nil || (tt.want != nil && !errors.Is(err, tt.want)) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
			if !reflect.DeepEqual(config, Config{}) {
				t.Fatalf("invalid input returned partial config: %+v", config)
			}
		})
	}
}

func TestConfigEnvironmentsAndBoundaries(t *testing.T) {
	for _, environment := range []string{"development", "testing", "production"} {
		for _, port := range []int{1, 65535} {
			config := Config{AppName: "app", Environment: environment, Port: port}
			data, err := EncodeConfig(config)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := DecodeConfig(data)
			if err != nil || !reflect.DeepEqual(decoded, config) {
				t.Fatalf("valid config round trip = %+v, %v", decoded, err)
			}
		}
	}
}

func TestEncodeConfigInvalid(t *testing.T) {
	for _, tt := range []struct {
		config Config
		want   error
	}{
		{Config{}, ErrEmptyAppName},
		{Config{AppName: "app", Environment: "testing", Port: 0}, ErrInvalidPort},
		{Config{AppName: "app", Environment: "other", Port: 8080}, ErrInvalidEnvironment},
	} {
		if data, err := EncodeConfig(tt.config); data != nil || !errors.Is(err, tt.want) {
			t.Fatalf("encode invalid config = %s, %v; want %v", data, err, tt.want)
		}
	}
}
