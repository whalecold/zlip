package encode

import (
	"fmt"
	"os"

	"whalecold/compress/pkg/entrance"

	"github.com/spf13/cobra"
)

const (
	encodeExample = `
# 解压文件
edm decode --source file --target file
`
)

type decodeOption struct {
	source string
	target string
}

// NewEncodeCommand return decode command
func NewEncodeCommand() *cobra.Command {
	opts := decodeOption{}
	cmd := &cobra.Command{
		Use:     "encode",
		Short:   "encode file.",
		Long:    "encode file.",
		Example: encodeExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.encode(args); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVar(&opts.source, "source", opts.source, "source file,【required】")
	cmd.Flags().StringVar(&opts.target, "target", opts.target, "target file,【required】")
	return cmd
}

func (d *decodeOption)encode(args []string) error {
	if len(d.target) == 0 || len(d.source) == 0 {
		return fmt.Errorf("source or target parmer can't be empty")
	}

	entrance.Entrance(d.source, d.target, true)
	return nil
}
