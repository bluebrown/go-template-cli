package main

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	flag "github.com/spf13/pflag"

	"github.com/Masterminds/sprig/v3"
	tm "github.com/bluebrown/treasure-map/pkg"
)

const rootTemplateName = "_tpl.root"

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// add a newline if the --no-newline flag was not set
	if !noNewline {
		fmt.Println()
	}
}

func run() error {
	var err error

	// create the root template
	tpl := template.New(rootTemplateName)
	tpl.Option(options...)
	tpl.Funcs(tm.MakeFuncMap(sprig.TxtFuncMap(), tpl))

	// parse the arguments
	for _, arg := range flag.Args() {
		tpl, err = tpl.Parse(arg)
		if err != nil {
			return err
		}
	}

	// parse files and globs in the order they were specified
	// to align with go's template package
	fileIndex := 0
	globIndex := 0
	for _, arg := range os.Args[1:] {
		if arg == "-f" || arg == "--file" {
			// parse next file
			file := files[fileIndex]
			tpl, err = tpl.ParseFiles(file)
			if err != nil {
				return err
			}
			fileIndex++
			continue
		}
		if arg == "-g" || arg == "--glob" {
			// parse next glob
			glob := globs[globIndex]
			tpl, err = tpl.ParseGlob(glob)
			if err != nil {
				return err
			}
			globIndex++
			continue
		}
	}

	// defined templates
	templates := tpl.Templates()

	// if there are no templates, return an error
	if len(templates) == 0 {
		return errors.New("no templates found")
	}

	// determine the template name to use
	if templateName == "" {
		if flag.NArg() > 0 {
			templateName = rootTemplateName
		} else if globIndex > 0 || fileIndex > 0 {
			templateName = templates[0].Name()
		}
	}

	// execute the template
	// read the input from stdin
	info, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	// data is used to store the decoded input
	var data any

	// if we are reading from stdin, decode the input
	if info.Mode()&os.ModeCharDevice == 0 {
		if fn, ok := decoderMap()[decoder]; ok {
			err = fn(os.Stdin, &data)
		} else {
			err = errors.New("unknown decoder")
		}
		if err != nil {
			return err
		}
	}

	// execute the template with the given name
	// and optional data from stdin
	return tpl.ExecuteTemplate(os.Stdout, templateName, data)
}
