package main

import (
	"fmt"
	"os"

	"github.com/vlnolam/vlnolam/orchestrator/internal/config"
)

func main() {
	stack, err := config.ParseFile("config.yaml")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed config: %+v\n", stack)
}
