package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

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
