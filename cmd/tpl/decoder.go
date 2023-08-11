package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/icza/dyno"
	"gopkg.in/yaml.v3"
)

var decoderMap = map[decoderkind]decode{
	decoderkindJSON: decodeJson,
	decoderkindYAML: decodeYaml,
	decoderkindTOML: decodeToml,
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
