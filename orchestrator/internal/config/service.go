package config

import (
	"strconv"

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
			err.Data["egress_index"] = strconv.Itoa(i)
			return err
		}
	}

	for i, ingress := range svc.Ingress {
		if err := ingress.Validate(IngressConnection); err != nil {
			err.Data["ingress_index"] = strconv.Itoa(i)
			return err
		}
	}

	return nil
}
