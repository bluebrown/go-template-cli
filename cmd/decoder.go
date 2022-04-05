package main

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/icza/dyno"
	"gopkg.in/yaml.v3"
)

func decoderMap() map[string]func(io.Reader, *any) error {
	return map[string]func(io.Reader, *any) error{
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
