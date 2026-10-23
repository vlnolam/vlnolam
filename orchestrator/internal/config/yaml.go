package config

import (
	"bytes"
	"fmt"

	"go.yaml.in/yaml/v4"
)

// decodeNode strictly decodes a node inside a custom UnmarshalYAML.
// A terrible workaround because Node.Decode hardcodes DefaultOptions (yaml issue #321).
func decodeStrict(n *yaml.Node, v any) error {
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(n); err != nil {
		return fmt.Errorf("encode node: %w", err)
	}

	dec := yaml.NewDecoder(&buf)

	// All this because of this line
	dec.KnownFields(true)

	if err := dec.Decode(v); err != nil {
		// Note: Line numbers are wrong. Ignore for now, wait for the yaml to provide a proper fix.
		return err
	}
	return nil
}
