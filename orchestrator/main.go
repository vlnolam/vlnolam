package main

import (
	"fmt"
	"os"

	"github.com/vlnolam/vlnolam/orchestrator/internal/config"
)

func main() {
	var stack config.Stack

	if err := stack.ParseFile("config.yaml"); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if err := stack.Validate(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed config: %+v\n", stack)
}
