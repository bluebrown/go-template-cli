package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

func helptext() {
	fmt.Print(`Usage of tpl:
  -t string
  -template string
	alternative way to specify template
  -n
  -no-newline
	do not print a new line at the end
  -h
  -help
	show this message
Examples:
Standard input:
	echo '{"place": "bar"}' | tpl 'lets go to the {{ .place }}!'
File:
	tpl -t path/to/template < path/to/input.json
`)
	os.Exit(0)
}

func main() {
	var t *template.Template
	var templatePath string
	var noNewline bool
	var templateStr = "{{.}}"
	var showHelp bool

	flag.StringVar(&templatePath, "t", "", "")
	flag.StringVar(&templatePath, "template", "", "alternative way to specify template")
	flag.BoolVar(&noNewline, "n", false, "")
	flag.BoolVar(&noNewline, "no-newline", false, "do not print a new line at the end")
	flag.BoolVar(&showHelp, "h", false, "")
	flag.BoolVar(&showHelp, "help", false, "show this message")
	flag.Parse()

	if showHelp {
		helptext()
		os.Exit(2)
	}

	if flag.NArg() > 1 {
		exit("too many arguments")
	}

	// try first to read template from file path
	if templatePath != "" {
		parts := strings.Split(templatePath, "/")
		name := parts[len(parts)-1]
		t = template.Must(baseTpl(name).ParseFiles(templatePath))
	} else {
		// if argument has been provided set the template
		if flag.NArg() == 1 {
			templateStr = flag.Arg(0)
		}
		// otherwise use default
		t = template.Must(baseTpl("arg").Parse(templateStr))
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
		exit(`Usage of tpl:
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
	}

	if err = t.Execute(os.Stdout, dyno.ConvertMapI2MapS(data)); err != nil {
		exit(err)
	}

	if !noNewline {
		fmt.Println()
	}
}

func baseTpl(name string) *template.Template {
	funcMap := map[string]interface{}{
		"toYaml": func(v interface{}) string {
			b, err := yaml.Marshal(v)
			if err != nil {
				return ""
			}
			return string(b)
		},
		"mustToYaml": func(v interface{}) (string, error) {
			b, err := yaml.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
	}
	return template.New(name).Funcs(funcMap).Funcs(sprig.TxtFuncMap())
}

func exit(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
