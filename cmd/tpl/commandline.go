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

// the options are used to control the runtime behavior of the commandLine
type commandLineOptions struct {
	// internal
	defaultTemplateName string
	// flags
	files        []string
	globs        []string
	templateName string
	options      []string
	decoder      decoder
	noNewline    bool
}

// generate options with default values
func defaultOptions() *commandLineOptions {
	return &commandLineOptions{
		defaultTemplateName: "_gotpl_default",
		decoder:             decodeJson,
	}
}

// link the opt struct fields to the given flagset. defaults to
// pflag.CommandLine. requires to call fs.Parse, AFTER calling initFlags
func (opts *commandLineOptions) initFlags(fs *pflag.FlagSet) {
	if fs == nil {
		fs = pflag.CommandLine
	}
	fs.StringArrayVarP(&opts.files, "file", "f", opts.files, "template file path. Can be specified multiple times")
	fs.StringArrayVarP(&opts.globs, "glob", "g", opts.globs, "template file glob. Can be specified multiple times")
	fs.StringVarP(&opts.templateName, "name", "n", opts.templateName, "if specified, execute the template with the given name")
	fs.VarP(&opts.decoder, "decoder", "d", "decoder to use for input data. Supported values: json, yaml, toml (default \"json\")")
	fs.StringArrayVar(&opts.options, "option", opts.options, "option to pass to the template engine. Can be specified multiple times")
	fs.BoolVar(&opts.noNewline, "no-newline", opts.noNewline, "do not print newline at the end of the output")
}

// run the template cli program
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
	tpl := baseTemplate(opts.defaultTemplateName, opts.options...)

	// parse positional arguments
	tpl, err = parsePosArgs(tpl, fs.Args())
	if err != nil {
		return fmt.Errorf("parse pos args: %w", err)
	}

	// parse optional arguments
	tpl, err = parseOptArgs(tpl, args, opts.files, opts.globs)
	if err != nil {
		return fmt.Errorf("parse opt args: %w", err)
	}

	// defined templates
	templates := tpl.Templates()

	// if there are no templates, return an error
	if len(templates) == 0 {
		return errors.New("no templates found")
	}

	// determine the template to use
	templateName, err := selectedTemplate(opts.templateName, opts.defaultTemplateName, fs.Args(), templates)
	if err != nil {
		return fmt.Errorf("select template: %w", err)
	}

	// data is used to store the decoded input
	var data any

	// decode the input stream, if any
	if input != nil {
		if err := opts.decoder(input, &data); err != nil {
			return fmt.Errorf("decode input: %w", err)
		}
	}

	// execute the template with the given name and optional data from stdin
	if err := tpl.ExecuteTemplate(output, templateName, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	// add a newline if the --no-newline flag was not set
	if !opts.noNewline {
		fmt.Fprintln(output)
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

// parse all positional arguments into the "default" template
func parsePosArgs(tpl *template.Template, posArgs []string) (*template.Template, error) {
	var err error
	for _, arg := range posArgs {
		tpl, err = tpl.Parse(arg)
		if err != nil {
			return nil, fmt.Errorf("error to parsing template: %v", err)
		}
	}
	return tpl, nil
}

// parse files and globs in the order they were specified, to align with go's
// template engine
func parseOptArgs(tpl *template.Template, rawArgs []string, files, globs []string) (*template.Template, error) {
	var (
		err       error
		fileIndex uint8
		globIndex uint8
	)

	// FIXME if arg is like --file=foo.txt, the if conditions wont detect it

	for _, arg := range rawArgs {
		if arg == "-f" || arg == "--file" {
			// parse next file
			file := files[fileIndex]
			tpl, err = tpl.ParseFiles(file)
			if err != nil {
				return nil, fmt.Errorf("error parsing file %s: %v", file, err)
			}
			fileIndex++
			continue
		}

		if arg == "-g" || arg == "--glob" {
			// parse next glob
			glob := globs[globIndex]
			tpl, err = tpl.ParseGlob(glob)
			if err != nil {
				return nil, fmt.Errorf("error parsing glob %s: %v", glob, err)
			}
			globIndex++
			continue
		}
	}

	return tpl, nil
}

// determine the template to execute. In the order of precedence:
//  1. current name, if set
//  2. default name, if at least 1 positional arg
//  3. templates name, if exactly 1 template
//  4. --name flag required, otherwise
func selectedTemplate(currentName, defaultName string, posArgs []string, templates []*template.Template) (string, error) {
	if currentName != "" {
		return currentName, nil
	}

	if len(posArgs) > 0 {
		return defaultName, nil
	}

	if len(templates) == 1 {
		return templates[0].Name(), nil
	}

	return "", fmt.Errorf("the --name flag is required when multiple templates are defined and no default template exists")
}
