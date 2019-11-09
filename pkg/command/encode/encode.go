package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/whalecold/zlip/pkg/entrance"
)

const (
	encodeExample = `
# 解压文件
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
		Use:     "encode",
		Short:   "encode file.",
		Long:    "encode file.",
		Example: encodeExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.encode(args); err != nil {
				panic(err)
			}
		},
	}
	cmd.Flags().StringVar(&opts.source, "source", opts.source, "source file,【required】")
	cmd.Flags().StringVar(&opts.target, "target", opts.target, "target file,【required】")
	return cmd
}

func (d *decodeOption) encode(args []string) error {
	if len(d.target) == 0 || len(d.source) == 0 {
		return fmt.Errorf("source or target parmer can't be empty")
	}

	entrance.Entrance(d.source, d.target, false)
	return nil
}
