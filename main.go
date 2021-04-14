package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

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

	if flag.NArg() == 1 {
		templateStr = flag.Arg(0)
		if templatePath != "" {
			exit("too many arguments")
		}
	}

	if templatePath != "" {
		content, err := ioutil.ReadFile(templatePath)
		if err != nil {
			exit(err)
		}
		templateStr = string(content)
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		exit(err)
	}

	var data interface{}
	if info.Mode()&os.ModeCharDevice == 0 {
		dec := json.NewDecoder(os.Stdin)
		for {
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				exit(err)
			}
		}
	}

	t, err := template.New("").Parse(`{{define "T"}}` + templateStr + `{{end}}`)
	if err != nil {
		exit(err)
	}
	if err = t.ExecuteTemplate(os.Stdout, "T", data); err != nil {
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
