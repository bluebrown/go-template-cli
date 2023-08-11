package main

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/spf13/pflag"
)

type state struct {
	// options
	defaultTemplateName string
	files               []string
	globs               []string
	templateName        string
	options             []string
	decoder             decoder
	noNewline           bool

	// internal state
	flagSet  *pflag.FlagSet
	template *template.Template
	rawArgs  []string
	posArgs  []string
}

func new(fs *pflag.FlagSet) *state {
	if fs == nil {
		fs = pflag.CommandLine
	}

	cli := &state{
		flagSet:             fs,
		decoder:             decodeJson,
		defaultTemplateName: "_gotpl_default",
	}

	fs.StringArrayVarP(&cli.files, "file", "f", cli.files, "template file path. Can be specified multiple times")
	fs.StringArrayVarP(&cli.globs, "glob", "g", cli.globs, "template file glob. Can be specified multiple times")
	fs.StringVarP(&cli.templateName, "name", "n", cli.templateName, "if specified, execute the template with the given name")
	fs.VarP(&cli.decoder, "decoder", "d", "decoder to use for input data. Supported values: json, yaml, toml (default \"json\")")
	fs.StringArrayVar(&cli.options, "option", cli.options, "option to pass to the template engine. Can be specified multiple times")
	fs.BoolVar(&cli.noNewline, "no-newline", cli.noNewline, "do not print newline at the end of the output")

	return cli
}

func (cli *state) run(args []string, r io.Reader, w io.Writer) (err error) {
	if err := cli.parse(args); err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	data, err := cli.decode(r)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	if err := cli.render(w, data); err != nil {
		return fmt.Errorf("render: %w", err)
	}

	return nil
}

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
	err := new(nil).run(os.Args[1:], input, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}
