package command

import (
	"io"

	"github.com/whalecold/zlip/pkg/command/decode"
	"github.com/whalecold/zlip/pkg/command/encode"
	"github.com/whalecold/zlip/pkg/version"

	"github.com/spf13/cobra"
)

// New returns compress command
func New(_ io.Reader, _, _ io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edm",
		Short: "encode or decode the file.",
		Long:  "encode or decode the file.",
	}

	cmd.AddCommand(decode.New())
	cmd.AddCommand(encode.New())
	cmd.AddCommand(NewVersionCommand())

	return cmd
}

// NewVersionCommand returns version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get version.",
		Long:  "Get version.",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintInfoAndExit()
		},
	}

	return cmd
}
