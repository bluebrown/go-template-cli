package textfunc

import (
	"bytes"
	"text/template"
)

func MakeInclude(t *template.Template) func(string, any) string {
	return func(name string, v any) string {
		b := new(bytes.Buffer)
		if err := t.ExecuteTemplate(b, name, v); err != nil {
			return ""
		}
		return b.String()
	}
}

func MakeTpl(t *template.Template) func(string, any) string {
	return func(snippet string, v any) string {
		t, err := t.Parse(snippet)
		if err != nil {
			return ""
		}
		b := new(bytes.Buffer)
		if err := t.Execute(b, v); err != nil {
			return ""
		}
		return b.String()
	}
}
