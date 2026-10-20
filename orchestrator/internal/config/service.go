package config

import "github.com/vlnolam/vlnolam/orchestrator/internal/errors"

type ServiceConfig struct {
	Name string `yaml:"name"`
}

func (sc *ServiceConfig) Validate() *errors.Error {
	if sc.Name == "" {
		return errors.New(errors.ErrConfigValidation, "service name is required")
	}
	return nil
}
