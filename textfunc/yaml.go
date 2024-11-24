package textfunc

import (
	"gopkg.in/yaml.v3"
)

func ToYaml(v any) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func MustToYaml(v any) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func FromYaml(s string) any {
	var v any
	err := yaml.Unmarshal([]byte(s), &v)
	if err != nil {
		return nil
	}
	return v
}
