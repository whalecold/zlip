package decode

import (
	"fmt"

	"github.com/whalecold/zlip/pkg/caller"
	"github.com/whalecold/zlip/pkg/caller/scheduler/processor"

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

// New return decode command
func New() *cobra.Command {
	opts := decodeOption{}
	cmd := &cobra.Command{
		Use:     "decode",
		Short:   "decode file.",
		Long:    "decode file.",
		Example: decodeExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.decode(args); err != nil {
				panic(err)
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

	caller.Run(d.source, d.target, processor.DecodeType)
	return nil
}
