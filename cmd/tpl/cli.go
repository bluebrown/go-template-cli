package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig/v3"
	"github.com/mlabbe/treasure-map/textfunc"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"os"
)

var version = "github.com/mlabbe/go-template-cli"

// always strict
var FatalMissingInclude = true
var TemplateOptions = "missingkey=error"

// the state of the program
type state struct {
	// options
	defaultTemplateName string
	files               []string
	globs               []string
	templateName        string
	decoder             decoder
	noNewline           bool
	showVersion         bool
	preservePreamble    bool
	outputFilename      string

	// internal state
	flagSet  *pflag.FlagSet
	template *template.Template
}

// create a new cli instance and bind flags to it
// flag.Parse is called on run
func new(fs *pflag.FlagSet) *state {
	if fs == nil {
		fs = pflag.CommandLine
	}

	cli := &state{
		flagSet:             fs,
		decoder:             decodeToml,
		defaultTemplateName: "_gotpl_default",
	}

	fs.StringArrayVarP(&cli.files, "file", "f", cli.files, "template file path. Can be specified multiple times")
	fs.StringArrayVarP(&cli.globs, "glob", "g", cli.globs, "template file glob. Can be specified multiple times")
	fs.StringVarP(&cli.templateName, "name", "n", cli.templateName, "if specified, execute the template with the given name")
	fs.StringVarP(&cli.outputFilename, "output-file", "o", "", "output filename (outputs to stdout if unspecified)")
	fs.VarP(&cli.decoder, "decoder", "d", "decoder to use for input data. Supported values: json, yaml, toml (default \"toml\")")
	fs.BoolVar(&cli.noNewline, "no-newline", cli.noNewline, "do not print newline at the end of the output")
	fs.BoolVar(&cli.showVersion, "version", cli.showVersion, "show version information and exit")
	fs.BoolVar(&cli.preservePreamble, "preserve-preamble", cli.preservePreamble, "Preserve build edge psecification comments in output file")

	return cli
}

func (cli *state) replaceOutputWriterFromCli(w io.Writer) (io.Writer, error) {

	if cli.outputFilename == "" {
		return w, nil
	}

	file, err := os.Create(cli.outputFilename)
	if err != nil {
		return nil, err
	}

	return file, err
}

// parse the options and input, decode the input and render the result
func (cli *state) run(args []string, r io.Reader, w io.Writer) (err error) {
	if err := cli.parse(args); err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	w, err = cli.replaceOutputWriterFromCli(w)
	if err != nil {
		return fmt.Errorf("output file: %w", err)
	}

	if cli.showVersion {
		fmt.Fprintln(w, version)
		return nil
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
	cli.flagSet.SortFlags = false

	if err := cli.flagSet.Parse(rawArgs); err != nil {
		return err
	}

	cli.template = baseTemplate(cli.defaultTemplateName)

	return nil
}

// parse all positional arguments into the "default" template. should be called
// after parseFlagset
func (cli *state) parsePositional() (err error) {
	for _, arg := range cli.flagSet.Args() {
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
	cli.flagSet.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "file":
			file := cli.files[fileIndex]
			cli.template, err = cli.template.ParseFiles(file)
			if err != nil {
				err = fmt.Errorf("error parsing file %s: %v", file, err)
				return
			}
			fileIndex++
		case "glob":
			glob := cli.globs[globIndex]
			cli.template, err = cli.template.ParseGlob(glob)
			if err != nil {
				err = fmt.Errorf("error parsing glob %s: %v", glob, err)
				return
			}
			globIndex++
		}
	})
	return cli.template, err
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

	if len(cli.flagSet.Args()) > 0 {
		return cli.defaultTemplateName, nil
	}

	if len(templates) == 1 {
		return templates[0].Name(), nil
	}

	return "", fmt.Errorf("the --name flag is required when multiple templates are defined and no default template exists")
}

// construct a base templates with custom functions attached
func baseTemplate(defaultName string) *template.Template {

	tpl := template.New(defaultName)
	tpl = tpl.Option(TemplateOptions)
	tpl = tpl.Funcs(textfunc.MapClosure(sprig.TxtFuncMap(), tpl, FatalMissingInclude))
	return tpl
}
