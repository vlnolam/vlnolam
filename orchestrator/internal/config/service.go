package config

import (
	"fmt"

	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
)

type Service struct {
	Name    string
	Egress  []Connection
	Ingress []Connection
}

func (svc *Service) WrapError(e *errors.Error) *errors.Error {
	e.Msg = fmt.Sprintf("service '%s': ", svc.Name) + e.Msg
	return e
}

func (svc *Service) Validate() *errors.Error {
	if svc.Name == "" {
		return errors.New(errors.ErrConfigValidation, "service name is required")
	}

	for _, egress := range svc.Egress {
		if err := egress.Validate(EgressConnection); err != nil {
			return err
		}
	}

	for _, ingress := range svc.Ingress {
		if err := ingress.Validate(IngressConnection); err != nil {
			return err
		}
	}

	return nil
}
