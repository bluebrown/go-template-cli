package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	files        []string
	globs        []string
	templateName string
	options      []string
	decoder      string
	noNewline    bool
	version      = "0.1.0"
	commit       = "unknown"
	debugMode    bool
	reflectMode  bool
)

func init() {
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

func init() {
	var showHelp bool
	var showVersion bool
	flag.StringArrayVarP(&files, "file", "f", []string{}, "template file path. Can be specified multiple times")
	flag.StringArrayVarP(&globs, "glob", "g", []string{}, "template file glob. Can be specified multiple times")
	flag.StringVarP(&templateName, "name", "n", "", "if specified, execute the template with the given name")
	flag.StringArrayVar(&options, "options", []string{}, "options to pass to the template engine")
	flag.StringVarP(&decoder, "decoder", "d", "json", "decoder to use for input data. Supported values: json, yaml, toml, xml")
	flag.BoolVar(&noNewline, "no-newline", false, "do not print newline at the end of the output")
	flag.BoolVarP(&showHelp, "help", "h", false, "show the help text")
	flag.BoolVarP(&showVersion, "version", "v", false, "show the version")
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if showVersion {
		fmt.Printf("tpl - version: %s, commit: %s\n", version, commit)
		os.Exit(0)
	}

	if debugMode {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

}
