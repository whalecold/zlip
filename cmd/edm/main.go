package main

import (
	goflag "flag"
	"fmt"
	"os"

	"whalecold/compress/pkg/command"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

func main() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Set("logtostderr", "true")
	defer glog.Flush()

	command := command.NewCompressAdmin(os.Stdin, os.Stdout, os.Stderr)
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
