package textfunc

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

func ToToml(v any) string {
	b := new(bytes.Buffer)
	enc := toml.NewEncoder(b)
	err := enc.Encode(v)
	if err != nil {
		return ""
	}
	return b.String()
}

func MustToToml(v any) (string, error) {
	b := new(bytes.Buffer)
	enc := toml.NewEncoder(b)
	err := enc.Encode(v)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func FromToml(s string) any {
	var v any
	_ = toml.Unmarshal([]byte(s), &v)
	return v
}
