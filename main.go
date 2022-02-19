package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

func main() {
	var t *template.Template
	var templatePath string
	var noNewline bool
	var templateStr = "{{.}}"

	flag.StringVar(&templatePath, "t", "", "")
	flag.StringVar(&templatePath, "template", "", "alternative way to specify template")
	flag.BoolVar(&noNewline, "n", false, "")
	flag.BoolVar(&noNewline, "no-newline", false, "do not print a new line at the end")
	flag.Parse()

	if flag.NArg() > 1 {
		exit("too many arguments")
	}

	// try first to read template from file path
	if templatePath != "" {
		parts := strings.Split(templatePath, "/")
		name := parts[len(parts)-1]
		t = template.Must(template.New(name).Funcs(sprig.TxtFuncMap()).ParseFiles(templatePath))
	} else {
		// if argument has been provided set the template
		if flag.NArg() == 1 {
			templateStr = flag.Arg(0)
		}
		// otherwise use default
		t = template.Must(template.New("any").Funcs(sprig.TxtFuncMap()).Parse(templateStr))
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		exit(err)
	}

	var data interface{}
	if info.Mode()&os.ModeCharDevice == 0 {
		// read the input data from stdin
		// must be valid json
		dec := yaml.NewDecoder(os.Stdin)
		for {
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				exit(err)
			}
		}
	}

	if data == nil {
		fmt.Fprintf(os.Stderr, `Usage of tpl:
  -n
  -no-newline
    do not print a new line at the end
  -t string
  -template string
    alternative way to specify template
Examples:
  Standard input:
    echo '{"place": "bar"}' | tpl 'lets go to the {{.place}}!'
  File:
    tpl -t path/to/template < path/to/input.json
ERROR: no data provided
`)
		os.Exit(1)
	}

	if err = t.Execute(os.Stdout, data); err != nil {
		exit(err)
	}

	if !noNewline {
		fmt.Println()
	}
}

func exit(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
