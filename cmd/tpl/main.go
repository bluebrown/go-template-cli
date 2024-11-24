package main

import (
	"fmt"
	"io"
	"os"

	cli "github.com/bluebrown/go-template-cli"
)

var version = "dev"

func main() {
	// reader may be used, if possible
	var input io.Reader

	// check if stdin is readable, and set the reader, if so
	if info, err := os.Stdin.Stat(); err == nil {
		if info.Mode()&os.ModeCharDevice == 0 {
			input = os.Stdin
		}
	}

	// try to run the program
	err := cli.New(version, nil).Run(os.Args[1:], input, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}
