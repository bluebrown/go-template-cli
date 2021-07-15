package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var t *template.Template

func main() {

	var templatePath string
	var newline bool
	var templateStr = "{{.}}"

	flag.StringVar(&templatePath, "t", "", "")
	flag.StringVar(&templatePath, "template", "", "alternative way to specify template")

	flag.BoolVar(&newline, "n", false, "")
	flag.BoolVar(&newline, "newline", false, "print new line at the end")

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
		dec := json.NewDecoder(os.Stdin)
		for {
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				exit(err)
			}
		}
	}

	if err = t.Execute(os.Stdout, data); err != nil {
		exit(err)
	}

	if newline {
		fmt.Println()
	}
}

func exit(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
