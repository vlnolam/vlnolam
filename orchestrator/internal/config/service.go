package config

import (
	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
)

type Service struct {
	Name    string
	Egress  []Connection
	Ingress []Connection
}

func (svc *Service) Validate() *errors.Error {
	if svc.Name == "" {
		return errors.New("service_name_required", nil)
	}

	for i, egress := range svc.Egress {
		if err := egress.Validate(EgressConnection); err != nil {
			err.PrependPath("egress/%d", i)
			return err
		}
	}

	for i, ingress := range svc.Ingress {
		if err := ingress.Validate(IngressConnection); err != nil {
			err.PrependPath("ingress/%d", i)
			return err
		}
	}

	return nil
}
