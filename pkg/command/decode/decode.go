package decode

import (
	"fmt"
	"os"

	"github.com/whalecold/zlip/pkg/entrance"

	"github.com/spf13/cobra"
)

const (
	decodeExample = `
# 压缩文件
zlip decode --source file --target file
`
)

type decodeOption struct {
	source string
	target string
}

// NewDecodeCommand return decode command
func NewDecodeCommand() *cobra.Command {
	opts := decodeOption{}
	cmd := &cobra.Command{
		Use:     "decode",
		Short:   "decode file.",
		Long:    "decode file.",
		Example: decodeExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.decode(args); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVar(&opts.source, "source", opts.source, "source file,【required】")
	cmd.Flags().StringVar(&opts.target, "target", opts.target, "target file,【required】")
	return cmd
}

func (d *decodeOption) decode(args []string) error {
	if len(d.target) == 0 || len(d.source) == 0 {
		return fmt.Errorf("source or target parmer can't be empty")
	}

	entrance.Entrance(d.source, d.target, true)
	return nil
}
