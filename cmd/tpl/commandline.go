package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/bluebrown/treasure-map/textfunc"
	"github.com/spf13/pflag"
)

type commandLineOptions struct {
	// internal
	defaultTemplateName string
	// flags
	files        []string
	globs        []string
	templateName string
	options      []string
	decoder      decoderkind
	noNewline    bool
}

func defaultOptions() *commandLineOptions {
	return &commandLineOptions{
		defaultTemplateName: "_gotpl_default",
		decoder:             decoderkindJSON,
	}
}

func (opts *commandLineOptions) initFlags(fs *pflag.FlagSet) {
	if fs == nil {
		fs = pflag.CommandLine
	}
	fs.StringArrayVarP(&opts.files, "file", "f", opts.files, "template file path. Can be specified multiple times")
	fs.StringArrayVarP(&opts.globs, "glob", "g", opts.globs, "template file glob. Can be specified multiple times")
	fs.StringVarP(&opts.templateName, "name", "n", opts.templateName, "if specified, execute the template with the given name")
	fs.VarP(&opts.decoder, "decoder", "d", "decoder to use for input data. Supported values: json, yaml, toml")
	fs.StringArrayVar(&opts.options, "option", opts.options, "option to pass to the template engine. Can be specified multiple times")
	fs.BoolVar(&opts.noNewline, "no-newline", opts.noNewline, "do not print newline at the end of the output")
}

func commandLine(ctx context.Context, args []string, input io.Reader, output io.Writer) (err error) {
	// parse the pflag set
	fs := pflag.NewFlagSet("tpl", pflag.ContinueOnError)
	fs.SetOutput(os.Stderr) // TODO handle the output

	opts := defaultOptions()
	opts.initFlags(fs)

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	// create the default template
	tpl := template.New(opts.defaultTemplateName)
	tpl.Option(opts.options...)
	tpl.Funcs(textfunc.MapClosure(sprig.TxtFuncMap(), tpl))

	// parse the arguments
	for _, arg := range fs.Args() {
		tpl, err = tpl.Parse(arg)
		if err != nil {
			return fmt.Errorf("error to parsing template: %v", err)
		}
	}

	// FIXME if arg is like --file=foo.txt,
	// the if conditions wont detect it
	//
	// parse files and globs in the order they were specified
	// to align with go's template package
	fileIndex := 0
	globIndex := 0
	for _, arg := range args {
		if arg == "-f" || arg == "--file" {
			// parse next file
			file := opts.files[fileIndex]
			tpl, err = tpl.ParseFiles(file)
			if err != nil {
				return fmt.Errorf("error parsing file %s: %v", file, err)
			}
			fileIndex++
			continue
		}
		if arg == "-g" || arg == "--glob" {
			// parse next glob
			glob := opts.globs[globIndex]
			tpl, err = tpl.ParseGlob(glob)
			if err != nil {
				return fmt.Errorf("error parsing glob %s: %v", glob, err)
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

	// determine the template to use
	if opts.templateName == "" {
		if len(fs.Args()) > 0 {
			opts.templateName = opts.defaultTemplateName
		} else if len(templates) == 1 {
			opts.templateName = templates[0].Name()
		} else {
			return fmt.Errorf(
				"the --name flag is required when multiple templates are defined and no default template exists%s",
				tpl.DefinedTemplates(),
			)
		}
	}

	// data is used to store the decoded input
	var data any

	// if there is a reader, decode the data from it
	if input != nil {
		if err := decoderMap[opts.decoder](input, &data); err != nil {
			return fmt.Errorf("error decoding input: %v", err)
		}
	}

	// execute the template with the given name
	// and optional data from stdin
	if err := tpl.ExecuteTemplate(output, opts.templateName, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	// add a newline if the --no-newline flag was not set
	if !opts.noNewline {
		fmt.Fprintln(output)
	}

	return nil
}
