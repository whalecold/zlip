package main

import (
	"fmt"
	"os"

	"github.com/whalecold/zlip/pkg/command"
)

func main() {
	command := command.NewCompressAdmin(os.Stdin, os.Stdout, os.Stderr)
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
