package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig/v3"
	"github.com/bluebrown/treasure-map/textfunc"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// the state of the program
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

// create a new cli instance and bind flags to it
// flag.Parse is called on run
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

// parse the options and input, decode the input and render the result
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

// decode the input stream into context data
func (cli *state) decode(r io.Reader) (any, error) {
	if r == nil || cli.decoder == nil {
		return nil, nil
	}
	var data any
	err := cli.decoder(r, &data)
	return data, err
}

// render a template
func (cli *state) render(w io.Writer, data any) error {
	templateName, err := cli.selectTemplate()
	if err != nil {
		return fmt.Errorf("select template: %w", err)
	}

	if err := cli.template.ExecuteTemplate(w, templateName, data); err != nil {
		return fmt.Errorf("execute template: %v", err)
	}

	if !cli.noNewline {
		fmt.Fprintln(w)
	}

	return nil
}

// determine the template to execute. In the order of precedence:
//  1. current name, if set
//  2. default name, if at least 1 positional arg
//  3. templates name, if exactly 1 template
//  4. --name flag required, otherwise
func (cli *state) selectTemplate() (string, error) {
	templates := cli.template.Templates()

	if len(templates) == 0 {
		return "", errors.New("no templates found")
	}

	if cli.templateName != "" {
		return cli.templateName, nil
	}

	if len(cli.posArgs) > 0 {
		return cli.defaultTemplateName, nil
	}

	if len(templates) == 1 {
		return templates[0].Name(), nil
	}

	return "", fmt.Errorf("the --name flag is required when multiple templates are defined and no default template exists")
}

type decoder func(in io.Reader, out any) error

func (dec decoder) String() string { return "" }

func (dec *decoder) Type() string { return "func" }

func (dec *decoder) Set(kind string) error {
	switch kind {
	case "json":
		*dec = decodeJson
	case "yaml":
		*dec = decodeYaml
	case "toml":
		*dec = decodeToml
	default:
		return fmt.Errorf("unsupported decoder %q", kind)
	}
	return nil
}

func decodeYaml(in io.Reader, out any) error {
	dec := yaml.NewDecoder(in)
	for {
		err := dec.Decode(out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func decodeToml(in io.Reader, out any) error {
	dec := toml.NewDecoder(in)
	_, err := dec.Decode(out)
	return err
}

func decodeJson(in io.Reader, out any) error {
	dec := json.NewDecoder(in)
	for {
		err := dec.Decode(out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
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
