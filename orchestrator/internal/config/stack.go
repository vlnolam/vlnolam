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
		return errors.New(errors.ErrConfigValidation, "stack name is required")
	}

	if sc.Version < minVersion || sc.Version > maxVersion {
		return errors.New(errors.ErrConfigValidation, "version must be between %d and %d", minVersion, maxVersion)
	}

	for _, svc := range sc.Services {
		if err := svc.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func ParseString(data string) (*Stack, *errors.Error) {
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(data)))
	decoder.KnownFields(true)

	var stack Stack
	if err := decoder.Decode(&stack); err != nil {
		return nil, errors.Wrap(errors.ErrConfigParsing, err)
	}

	if err := stack.Validate(); err != nil {
		return nil, err
	}

	return &stack, nil
}

func ParseFile(path string) (*Stack, *errors.Error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(errors.ErrFile, err)
	}

	return ParseString(string(data))
}
