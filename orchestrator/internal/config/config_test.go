package config

import (
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

			var scenarioMap map[string]interface{}
			if err := yaml.Unmarshal(scenarioData, &scenarioMap); err != nil {
				t.Fatalf("Failed to parse scenario file: %v", err)
			}

			configYaml, _ := yaml.Marshal(scenarioMap["config"])

			stack := Stack{}
			err := stack.ParseBytes(configYaml)
			if err == nil {
				err = stack.Validate()
			}

			expectedError := scenarioMap["expected_error"].(string)

			if err == nil {
				t.Errorf("Expected error %q but parsing succeeded", expectedError)
				return
			}

			if err.Data["key"] != expectedError {
				t.Errorf("Expected error %q but got %q", expectedError, err.Data["key"])
				t.Logf("Full error: %+v", err)
				return
			}

			t.Log("OK")
		})
	}
}
