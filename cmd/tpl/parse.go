package main

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/bluebrown/treasure-map/textfunc"
)

// parse options and args / files into templates. this method builds up the
// state of the program. It revolves primarily around the parsed templates. The
// args are parsed in a specific order:
//  1. parse the options into the flagset
//  2. parse positional arguments into the default template
//  3. parse files and globs in the order specified
func (cli *state) parse(rawArgs []string) error {
	if err := cli.parseFlagset(rawArgs); err != nil {
		return fmt.Errorf("parse raw args: %s", err)
	}

	if err := cli.parsePositional(); err != nil {
		return fmt.Errorf("parse pos args: %w", err)
	}

	if _, err := cli.parseFilesAndGlobs(); err != nil {
		return fmt.Errorf("parse opt args: %w", err)
	}

	return nil
}

// construct a base templates with custom functions attached
func baseTemplate(defaultName string, options ...string) *template.Template {
	tpl := template.New(defaultName)
	tpl = tpl.Option(options...)
	tpl = tpl.Funcs(textfunc.MapClosure(sprig.TxtFuncMap(), tpl))
	return tpl
}

func (cli *state) parseFlagset(rawArgs []string) error {
	if err := cli.flagSet.Parse(rawArgs); err != nil {
		return err
	}

	cli.rawArgs = rawArgs
	cli.posArgs = cli.flagSet.Args()

	cli.template = baseTemplate(cli.defaultTemplateName, cli.options...)

	return nil
}

// parse all positional arguments into the "default" template. should be called
// after parseFlagset
func (cli *state) parsePositional() (err error) {
	for _, arg := range cli.posArgs {
		cli.template, err = cli.template.Parse(arg)
		if err != nil {
			return fmt.Errorf("parse template: %v", err)
		}
	}
	return nil
}

// parse files and globs in the order they were specified, to align with go's
// template engine. should be called after parseFlagset
func (cli *state) parseFilesAndGlobs() (*template.Template, error) {
	var (
		err       error
		fileIndex uint8
		globIndex uint8
	)

	// FIXME if arg is like --file=foo.txt,
	// the if conditions wont detect it

	for _, arg := range cli.rawArgs {
		if arg == "-f" || arg == "--file" {
			// parse next file
			file := cli.files[fileIndex]
			cli.template, err = cli.template.ParseFiles(file)
			if err != nil {
				return nil, fmt.Errorf("error parsing file %s: %v", file, err)
			}
			fileIndex++
			continue
		}

		if arg == "-g" || arg == "--glob" {
			// parse next glob
			glob := cli.globs[globIndex]
			cli.template, err = cli.template.ParseGlob(glob)
			if err != nil {
				return nil, fmt.Errorf("error parsing glob %s: %v", glob, err)
			}
			globIndex++
			continue
		}
	}

	return cli.template, nil
}
