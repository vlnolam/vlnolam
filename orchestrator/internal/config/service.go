package config

import "github.com/vlnolam/vlnolam/orchestrator/internal/errors"

type Service struct {
	Name    string
	Egress  []Connection
	Ingress []Connection
}

func (sc *Service) Validate() *errors.Error {
	if sc.Name == "" {
		return errors.New(errors.ErrConfigValidation, "service name is required")
	}
	return nil
}
