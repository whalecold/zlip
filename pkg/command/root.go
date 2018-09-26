package command

import (
	"io"

	"whalecold/compress/pkg/command/decode"
	"whalecold/compress/pkg/command/encode"
	"whalecold/compress/pkg/version"

	"github.com/spf13/cobra"
)

// NewCompressAdmin returns compress command
func NewCompressAdmin(_ io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "edm",
		Short: "encode or decode the file.",
		Long:  "encode or decode the file.",
	}

	cmds.AddCommand(decode.NewDecodeCommand())
	cmds.AddCommand(encode.NewEncodeCommand())
	cmds.AddCommand(NewVersionCommand())

	return cmds
}

// NewVersionCommand returns version command
func NewVersionCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "version",
		Short: "Get version.",
		Long:  "Get version.",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintInfoAndExit()
		},
	}

	return cmds
}