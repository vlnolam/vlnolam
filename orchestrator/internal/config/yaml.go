package config

import (
	"bytes"

	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
	"go.yaml.in/yaml/v4"
)

func wrapParseError(error error) *errors.Error {
	if err, ok := error.(*errors.Error); ok {
		// Already our error, keep as is
		return err
	}

	if err, ok := error.(*yaml.LoadErrors); ok && len(err.Errors) > 0 {
		return wrapYamlError(err.Errors[0])
	}

	if err, ok := error.(*yaml.LoadError); ok {
		return wrapYamlError(err)
	}

	return errors.Wrap("config_parse_error", error)
}

func wrapYamlError(error *yaml.LoadError) *errors.Error {
	var result *errors.Error

	if err, ok := error.Unwrap().(*errors.Error); ok {
		result = err
	} else {
		result = errors.Wrap("config_parse_error", error)
	}

	if line := error.Mark.Line; line > 0 {
		result.Data["line"] = line
	}

	return result
}

// decodeNode strictly decodes a node inside a custom UnmarshalYAML.
// A terrible workaround because Node.Decode hardcodes DefaultOptions (yaml issue #321).
func decodeStrict(n *yaml.Node, v any) error {
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(n); err != nil {
		return err
	}

	dec := yaml.NewDecoder(&buf)

	// All this because of this line
	dec.KnownFields(true)

	if err := dec.Decode(v); err != nil {
		// Note: Line numbers are wrong. Ignore for now, wait for the yaml to provide a proper fix.
		return wrapParseError(err)
	}

	return nil
}
