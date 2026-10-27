package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vlnolam/vlnolam/orchestrator/internal/config"
	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
	"github.com/vlnolam/vlnolam/orchestrator/internal/format"
)

type GlobalParams struct {
	Formatter format.Formatter
}

func Execute() {
	var params GlobalParams
	params.Formatter = &format.BubbleFormatter{}

	if err := newCLI(&params).Execute(); err != nil {
		params.Formatter.Error(errors.Wrap("cli_error", err))
		params.Formatter.Finalize()
		os.Exit(1)
	}

	params.Formatter.Done()
	params.Formatter.Finalize()
}

func newCLI(params *GlobalParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "vlnolam",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newStackCommand(params))

	return cmd
}

type StackParams struct {
	*GlobalParams
	StackFile string
}

func newStackCommand(parentParams *GlobalParams) *cobra.Command {
	var params StackParams
	params.GlobalParams = parentParams

	cmd := &cobra.Command{
		Use:   "stack",
		Short: "Manage stack",
		Args:  cobra.NoArgs,
	}

	cmd.PersistentFlags().StringVarP(&params.StackFile, "file", "f", "vlnolam.yaml", "stack definition file")

	cmd.AddCommand(newStackUpCommand(&params))

	return cmd
}

func newStackUpCommand(parentParams *StackParams) *cobra.Command {
	var params struct {
		*StackParams
	}
	params.StackParams = parentParams

	result := &cobra.Command{
		Use:   "up",
		Short: "Validate the stack and start all services",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var stack config.Stack
			if err := loadStack(&stack, params.StackFile); err != nil {
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
