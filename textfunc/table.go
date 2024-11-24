package textfunc

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

func Table(v any) string {
	s, ok := v.([]any)
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
	case map[string]any:
		err = mapTable(w, s)
	case []any:
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

func mapTable(w *tabwriter.Writer, v []any) error {

	keys := []string{}
	var vals []string

	for i, row := range v {
		r, ok := row.(map[string]any)
		if !ok {
			return errors.New("mapTable: row is not map[string]any")
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

func sliceTable(w *tabwriter.Writer, v []any) error {

	var vals []string

	for i, row := range v {
		r, ok := row.([]any)
		if !ok {
			return errors.New("sliceTable: row is not []any")
		}

		if i == 0 {
			vals = make([]string, len(r))
		}

		for j, v := range r {
			if i == 0 {
				vals[j] = strings.Title(fmt.Sprintf("%v", v))
			} else {
				vals[j] = fmt.Sprintf("%v", v)
			}
		}

		if i == len(v)-1 {
			fmt.Fprint(w, strings.Join(vals, "\t"))
		} else {
			fmt.Fprintln(w, strings.Join(vals, "\t"))
		}
	}
	return nil
}
