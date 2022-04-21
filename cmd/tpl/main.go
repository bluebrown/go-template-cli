package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/icza/dyno"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/Masterminds/sprig/v3"
	"github.com/bluebrown/treasure-map/textfunc"
)

var (
	version = "0.2.0"
	commit  = "unknown"
)

var (
	files        []string
	globs        []string
	templateName string
	options      []string
	decoder      decoderkind = decoderkindJSON
	noNewline    bool
)

const defaultTemplateName = "_gotpl_default"

var decoderMap = map[decoderkind]decode{
	decoderkindJSON: decodeJson,
	decoderkindYAML: decodeYaml,
	decoderkindTOML: decodeToml,
}

func setFlagUsage() {
	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage: tpl [--file PATH]... [--glob PATTERN]... [--name TEMPLATE_NAME]
           [--decoder DECODER_NAME] [--option KEY=VALUE]... [--no-newline] [TEMPLATE...] [-]
           [--help] [--usage] [--version]`)
	}
}

func helptext() {
	fmt.Fprintln(os.Stderr, "Usage: tpl [options] [templates]")
	fmt.Fprintln(os.Stderr, "Options:")
	pflag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Examples")
	fmt.Fprintln(os.Stderr, "  tpl '{{ . }}' < data.json")
	fmt.Fprintln(os.Stderr, "  tpl --file my-template.tpl < data.json")
	fmt.Fprintln(os.Stderr, "  tpl --glob 'templates/*' --name foo.tpl < data.json")
}

func parseFlags() {
	var showHelp bool
	var showUsage bool
	var showVersion bool
	pflag.StringArrayVarP(&files, "file", "f", []string{}, "template file path. Can be specified multiple times")
	pflag.StringArrayVarP(&globs, "glob", "g", []string{}, "template file glob. Can be specified multiple times")
	pflag.StringVarP(&templateName, "name", "n", "", "if specified, execute the template with the given name")
	pflag.VarP(&decoder, "decoder", "d", "decoder to use for input data. Supported values: json, yaml, toml")
	pflag.StringArrayVar(&options, "option", []string{}, "option to pass to the template engine. Can be specified multiple times")
	pflag.BoolVar(&noNewline, "no-newline", false, "do not print newline at the end of the output")
	pflag.BoolVarP(&showHelp, "help", "h", false, "show the help text")
	pflag.BoolVar(&showUsage, "usage", false, "show the short usage text")
	pflag.BoolVarP(&showVersion, "version", "v", false, "show the version")
	pflag.Parse()
	if showHelp {
		helptext()
		os.Exit(0)
	}
	if showUsage {
		pflag.Usage()
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(os.Stderr, "version %s - commit %s\n", version, commit)
		os.Exit(0)
	}
}

func main() {
	setFlagUsage()
	parseFlags()

	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		pflag.Usage()
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(2)
	}

	// add a newline if the --no-newline flag was not set
	if !noNewline {
		fmt.Println()
	}
}

func run() (err error) {
	// create the root template
	tpl := template.New(defaultTemplateName)
	tpl.Option(options...)
	tpl.Funcs(textfunc.MapClosure(sprig.TxtFuncMap(), tpl))

	// parse the arguments
	for _, arg := range pflag.Args() {
		tpl, err = tpl.Parse(arg)
		if err != nil {
			return fmt.Errorf("error to parsing template: %v", err)
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
				return fmt.Errorf("error parsing file %s: %v", file, err)
			}
			fileIndex++
			continue
		}
		if arg == "-g" || arg == "--glob" {
			// parse next glob
			glob := globs[globIndex]
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
	if templateName == "" {
		if len(pflag.Args()) > 0 {
			templateName = defaultTemplateName
		} else if len(templates) == 1 {
			templateName = templates[0].Name()
		} else {
			return errors.New(fmt.Sprintf(
				"the --name flag is required when multiple templates are defined and no default template exists%s",
				tpl.DefinedTemplates(),
			))
		}
	}

	// execute the template
	// read the input from stdin
	info, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("error reading stdin: %v", err)
	}

	// data is used to store the decoded input
	var data any

	// if we are reading from stdin, decode the input
	if info.Mode()&os.ModeCharDevice == 0 {
		if err := decoderMap[decoder](os.Stdin, &data); err != nil {
			return fmt.Errorf("error decoding input: %v", err)
		}
	}

	// execute the template with the given name
	// and optional data from stdin
	if err := tpl.ExecuteTemplate(os.Stdout, templateName, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}

type decoderkind string // json, yaml, toml

const (
	decoderkindJSON decoderkind = "json"
	decoderkindYAML decoderkind = "yaml"
	decoderkindTOML decoderkind = "toml"
)

func (d *decoderkind) Set(s string) error {
	switch s {
	case "json", "yaml", "toml":
		*d = decoderkind(s)
		return nil
	default:
		return fmt.Errorf(
			"invalid decoder kind: %s, supported value are: %s, %s, %s",
			s,
			decoderkindJSON,
			decoderkindYAML,
			decoderkindTOML,
		)
	}
}

func (d *decoderkind) String() string {
	return string(*d)
}

func (d *decoderkind) Type() string {
	return "string"
}

type decode func(io.Reader, *any) error

func decodeYaml(in io.Reader, out *any) error {
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
	*out = dyno.ConvertMapI2MapS(*out)
	return nil
}

func decodeToml(in io.Reader, out *any) error {
	dec := toml.NewDecoder(in)
	_, err := dec.Decode(out)
	return err
}

func decodeJson(in io.Reader, out *any) error {
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
