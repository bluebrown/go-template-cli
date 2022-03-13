package funcs

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

func ToYaml(v interface{}) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func MustToYaml(v interface{}) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Table(v interface{}) string {
	s, ok := v.([]interface{})
	if !ok {
		return ""
	}

	if len(s) == 0 {
		return ""
	}

	w := new(tabwriter.Writer)
	buf := new(bytes.Buffer)
	w.Init(buf, 0, 8, 0, '\t', 0)
	var err error

	switch s[0].(type) {
	case map[string]interface{}:
		err = mapTable(w, s)
	case []interface{}:
		err = sliceTable(w, s)
	default:
		return ""
	}

	if err != nil {
		return ""
	}

	w.Flush()
	return buf.String()

}

func mapTable(w *tabwriter.Writer, v []interface{}) error {

	keys := []string{}
	var vals []string

	for i, row := range v {
		r, ok := row.(map[string]interface{})
		if !ok {
			return errors.New("mapTable: row is not map[string]interface{}")
		}

		if i == 0 {
			j := 0
			for k := range r {
				keys = append(keys, k)
				j++
			}
			vals = make([]string, len(keys))
			sort.Strings(keys)
			upperKeys := make([]string, len(keys))
			for i, k := range keys {
				upperKeys[i] = strings.ToUpper(k)
			}
			fmt.Fprintln(w, strings.Join(upperKeys, "\t"))
		}

		j := 0
		for _, k := range keys {
			v, ok := r[k]
			if !ok {
				vals[j] = ""
			} else {
				vals[j] = fmt.Sprintf("%v", v)
			}
			j++
		}

		if i == len(v)-1 {
			fmt.Fprint(w, strings.Join(vals, "\t"))
		} else {
			fmt.Fprintln(w, strings.Join(vals, "\t"))
		}
	}
	return nil
}

func sliceTable(w *tabwriter.Writer, v []interface{}) error {

	var vals []string

	for i, row := range v {
		r, ok := row.([]interface{})
		if !ok {
			return errors.New("sliceTable: row is not []interface{}")
		}

		if i == 0 {
			vals = make([]string, len(r))
		}

		for j, v := range r {
			vals[j] = fmt.Sprintf("%v", v)
		}

		if i == len(v)-1 {
			fmt.Fprint(w, strings.Join(vals, "\t"))
		} else {
			fmt.Fprintln(w, strings.Join(vals, "\t"))
		}
	}
	return nil
}
