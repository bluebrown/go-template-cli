package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/icza/dyno"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/Masterminds/sprig/v3"
	tm "github.com/bluebrown/treasure-map/pkg"
)

var (
	files        []string
	globs        []string
	templateName string
	options      []string
	decoder      string
	noNewline    bool
	version      = "0.1.1"
	commit       = "unknown"
	debugMode    bool
	reflectMode  bool
)

const rootTemplateName = "_tpl.root"

func main() {
	setFlagUsage()
	parseFlags()

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

type decode func(io.Reader, *any) error

// a map of decoder functions
func decoderMap() map[string]decode {
	return map[string]decode{
		"yaml": decodeYaml,
		"toml": decodeToml,
		"xml":  decodeXml,
		"json": decodeJson,
	}
}

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

func decodeXml(in io.Reader, out *any) error {
	dec := xml.NewDecoder(in)
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

func setFlagUsage() {
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] [templates]\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("Examples:\n")
		fmt.Printf("  tpl '{{ . }}' < data.json\n")
		fmt.Printf("  tpl --file my-template.tpl < data.json\n")
		fmt.Printf("  tpl --glob templates/* < data.json\n")
	}

}

func parseFlags() {
	var showHelp bool
	var showVersion bool
	flag.StringArrayVarP(&files, "file", "f", []string{}, "template file path. Can be specified multiple times")
	flag.StringArrayVarP(&globs, "glob", "g", []string{}, "template file glob. Can be specified multiple times")
	flag.StringVarP(&templateName, "name", "n", "", "if specified, execute the template with the given name")
	flag.StringVarP(&decoder, "decoder", "d", "json", "decoder to use for input data. Supported values: json, yaml, toml, xml")
	flag.StringArrayVar(&options, "options", []string{}, "options to pass to the template engine")
	flag.BoolVar(&noNewline, "no-newline", false, "do not print newline at the end of the output")
	flag.BoolVarP(&showHelp, "help", "h", false, "show the help text")
	flag.BoolVarP(&showVersion, "version", "v", false, "show the version")
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if showVersion {
		fmt.Printf("version %s - commit %s\n", version, commit)
		os.Exit(0)
	}

	if debugMode {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

}
