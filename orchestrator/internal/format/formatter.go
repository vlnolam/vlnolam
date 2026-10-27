package format

import "github.com/vlnolam/vlnolam/orchestrator/internal/errors"

type Formatter interface {
	Done()
	Error(*errors.Error)

	Finalize()
}
