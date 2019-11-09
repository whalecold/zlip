package main

import (
	"os"

	"github.com/whalecold/zlip/pkg/command"
)

func main() {
	cmd := command.New(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
