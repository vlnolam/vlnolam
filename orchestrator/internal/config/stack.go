package config

import (
	"bytes"
	"os"

	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
	"go.yaml.in/yaml/v4"
)

const minVersion = 1
const maxVersion = 1

type Stack struct {
	Version  int
	Stack    string
	Services []Service
}

func (sc *Stack) Validate() *errors.Error {
	if sc.Stack == "" {
		return errors.New("stack_name_required", nil)
	}

	if sc.Version < minVersion || sc.Version > maxVersion {
		return errors.New("unsupported_version", errors.Data{"min": maxVersion, "max": maxVersion})
	}

	existingServices := map[string]struct{}{}
	for i, svc := range sc.Services {
		if err := svc.Validate(); err != nil {
			if svc.Name != "" {
				err.PrependPath("services/%s", svc.Name)
			} else {
				err.PrependPath("services/%d", i)
			}
			return err
		}

		if _, found := existingServices[svc.Name]; found {
			return errors.New("duplicate_service_definition", errors.Data{"name": svc.Name})
		}

		existingServices[svc.Name] = struct{}{}
	}

	return nil
}

func ParseBytes(data []byte) (*Stack, *errors.Error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var stack Stack
	if err := decoder.Decode(&stack); err != nil {
		return nil, errors.Wrap("config_parse_error", err)
	}

	if err := stack.Validate(); err != nil {
		return nil, err
	}

	return &stack, nil
}

func ParseFile(path string) (*Stack, *errors.Error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap("file_error", err)
	}

	return ParseBytes(data)
}
