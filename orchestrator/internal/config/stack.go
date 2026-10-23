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
	Stack    string
	Services []Service
}

func (sc *Stack) Validate() *errors.Error {
	if sc.Stack == "" {
		return errors.New("stack_name_required", nil)
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

func (stack *Stack) ParseBytes(data []byte) *errors.Error {
	type versionHeader struct {
		Version int
	}

	// First validate the version
	// Otherwise we could get parse errors before checking the version value
	{
		decoder := yaml.NewDecoder(bytes.NewReader(data))
		// NO decoder.KnownFields(true)

		var header versionHeader
		if err := decoder.Decode(&header); err != nil {
			return errors.Wrap("config_parse_error", err)
		}

		if header.Version < minVersion || header.Version > maxVersion {
			return errors.New("unsupported_version", errors.Data{"min": minVersion, "max": maxVersion})
		}
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	type stackAndVersion struct {
		Header versionHeader `yaml:",inline"`
		Stack  Stack         `yaml:",inline"`
	}

	var stackExt stackAndVersion
	if err := decoder.Decode(&stackExt); err != nil {
		return errors.Wrap("config_parse_error", err)
	}

	*stack = stackExt.Stack
	return nil
}

func (stack *Stack) ParseFile(path string) *errors.Error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap("file_error", err)
	}

	return stack.ParseBytes(data)
}
