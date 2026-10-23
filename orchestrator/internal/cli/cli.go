package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vlnolam/vlnolam/orchestrator/internal/config"
	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
)

type GlobalArgs struct {
	JsonMode bool
}

func Execute() {
	var globalArgs GlobalArgs

	if err := newCLI(&globalArgs).Execute(); err != nil {
		errorData := errors.Wrap("cli_error", err).Data

		if globalArgs.JsonMode {
			enc := json.NewEncoder(os.Stderr)
			enc.Encode(map[string]any{"status": "error", "error": errorData})

		} else {
			fmt.Fprintf(os.Stderr, "[ERROR] %s\n", errorData["key"])
			for key, value := range errorData {
				if key == "key" || key == "message" {
					continue
				}

				fmt.Fprintf(os.Stderr, "  %s: %s\n", key, value)
			}

			if msg := errorData["message"]; msg != "" {
				fmt.Fprintf(os.Stderr, "\n%s\n", msg)
			}
		}
		os.Exit(-1)
	}

	if globalArgs.JsonMode {
		enc := json.NewEncoder(os.Stderr)
		enc.Encode(map[string]any{"status": "done"})
	}
}

func newCLI(args *GlobalArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vlnolam",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().BoolVar(&args.JsonMode, "json", false, "report in the JSON format")

	cmd.AddCommand(newStackCommand(args))

	return cmd
}

type StackArgs struct {
	*GlobalArgs
	StackFile string
}

func newStackCommand(parentArgs *GlobalArgs) *cobra.Command {
	var args StackArgs
	args.GlobalArgs = parentArgs

	cmd := &cobra.Command{
		Use:   "stack",
		Short: "Manage stack",
		Args:  cobra.NoArgs,
	}

	cmd.PersistentFlags().StringVarP(&args.StackFile, "file", "f", "vlnolam.yaml", "stack definition file")

	cmd.AddCommand(newStackUpCommand(&args))

	return cmd
}

func newStackUpCommand(parentArgs *StackArgs) *cobra.Command {
	var args struct {
		*StackArgs
	}
	args.StackArgs = parentArgs

	result := &cobra.Command{
		Use:   "up",
		Short: "Validate the stack and start all services",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, cmdArgs []string) error {
			var stack config.Stack
			if err := loadStack(&stack, args.StackFile); err != nil {
				return err
			}
			return nil
		},
	}

	return result
}

func loadStack(stack *config.Stack, file string) *errors.Error {
	if err := stack.ParseFile(file); err != nil {
		return err
	}
	if err := stack.Validate(); err != nil {
		return err
	}

	return nil
}
