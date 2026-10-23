package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestConfigErrors(t *testing.T) {
	scenarioDir := filepath.Join("testdata", "errors")

	entries, err := os.ReadDir(scenarioDir)
	if err != nil {
		t.Fatalf("Failed to read scenario directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			scenarioPath := filepath.Join(scenarioDir, entry.Name())
			scenarioData, _ := os.ReadFile(scenarioPath)

			var scenario struct {
				Config        yaml.Node
				ExpectedError string `yaml:"expected_error"`
			}

			decoder := yaml.NewDecoder(bytes.NewReader(scenarioData))
			decoder.KnownFields(true)
			decoder.Decode(&scenario)

			stack := Stack{}
			configYaml, _ := yaml.Marshal(scenario.Config) // node re-emits in original order, original style
			err := stack.ParseBytes(configYaml)
			if err == nil {
				err = stack.Validate()
			}

			if err == nil {
				t.Errorf("Expected error %q but parsing succeeded", scenario.ExpectedError)
				return
			}

			if err.Data["key"] != scenario.ExpectedError {
				t.Errorf("Expected error %q but got %q", scenario.ExpectedError, err.Data["key"])
				t.Logf("Full error: %+v", err)
				return
			}

			t.Log("OK")
		})
	}
}
