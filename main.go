package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/bluebrown/tpl/funcs"
	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

func main() {
	flag.Usage = func() {
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
`,
		)
	}

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
		flag.Usage()
		exit(0)
	}

	if flag.NArg() > 1 {
		exit("too many arguments")
	}

	// try first to read template from file path
	if templatePath != "" {
		t = template.Must(baseTpl(filepath.Base(templatePath)).ParseFiles(templatePath))
	} else {
		// if argument has been provided set the template
		if flag.NArg() == 1 {
			templateStr = flag.Arg(0)
		}
		// otherwise use default
		t = template.Must(baseTpl("positional.arg").Parse(templateStr))
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		exit(err)
	}

	var data interface{}
	if info.Mode()&os.ModeCharDevice == 0 {
		// read the input data from stdin
		// must be valid yaml/json
		dec := yaml.NewDecoder(os.Stdin)
		for {
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				exit(err)
			}
		}
	}

	if err = t.Execute(os.Stdout, dyno.ConvertMapI2MapS(data)); err != nil {
		exit(err)
	}

	if !noNewline {
		fmt.Println()
	}
}

func baseTpl(name string) *template.Template {
	return template.New(name).Funcs(funcs.TxtFuncMap()).Funcs(sprig.TxtFuncMap())
}

func exit(a ...interface{}) {
	flag.Usage()
	fmt.Print("Error: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
