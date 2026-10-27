package format

import (
	"fmt"
	"os"

	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
)

type SimpleFormatter struct{}

func (f *SimpleFormatter) Done() {
	fmt.Fprintln(os.Stderr, "[OK] Done")
}

func (f *SimpleFormatter) Finalize() {

}

func (f *SimpleFormatter) Error(e *errors.Error) {
	fmt.Fprintf(os.Stderr, "[ERROR] %s\n", e.Data["key"])
	for key, value := range e.Data {
		if key == "key" {
			continue
		}
		fmt.Fprintf(os.Stderr, "  %s: %s\n", key, value)
	}
}
