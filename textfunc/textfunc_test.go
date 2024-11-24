package textfunc_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/bluebrown/go-template-cli/textfunc"
)

func TestMapClosure(t *testing.T) {
	tpl := template.New("test")
	baseMap := make(map[string]any)
	fm := textfunc.MapClosure(baseMap, tpl)
	if _, ok := fm["include"]; !ok {
		t.Error("include is not found")
	}
}

func TestFuncMap(t *testing.T) {
	fm := textfunc.Map()
	if _, ok := fm["toYaml"]; !ok {
		t.Error("include is not found")
	}
}

func TestFuncs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		giveTemplate string
		giveData     any
		wantRendered string
	}{
		{
			name:         "tpl",
			giveTemplate: `{{tpl "{{ .foo.bar }}" . }}`,
			giveData:     map[string]any{"foo": map[string]any{"bar": "baz"}},
			wantRendered: "baz",
		},
		{
			name:         "include",
			giveTemplate: `{{ define "my.tpl" }}{{ .foo }}{{ end }}{{ include "my.tpl" . }}`,
			giveData:     map[string]any{"foo": "bar"},
			wantRendered: "bar",
		},
		{
			name:         "mapTable",
			giveTemplate: `{{ table . }}`,
			giveData:     []any{map[string]any{"a": 1, "b": 2}, map[string]any{"a": 3, "b": 4}},
			wantRendered: "A\tB\n1\t2\n3\t4",
		},
		{
			name:         "sliceTable",
			giveTemplate: `{{ table . }}`,
			giveData:     []any{[]any{"a", "b"}, []any{1, 2}, []any{3, 4}},
			wantRendered: "A\tB\n1\t2\n3\t4",
		},
		{
			name:         "toYaml",
			giveTemplate: `{{ toYaml . }}`,
			giveData:     map[string]any{"a": 1, "b": 2},
			wantRendered: "a: 1\nb: 2\n",
		},
		{
			name:         "mustToYaml",
			giveTemplate: `{{ mustToYaml . }}`,
			giveData:     map[string]any{"a": 1, "b": 2},
			wantRendered: "a: 1\nb: 2\n",
		},
		{
			name:         "fromYaml",
			giveTemplate: `{{ fromYaml .yaml }}`,
			giveData:     map[string]any{"yaml": "a: 1\nb: 2\n"},
			wantRendered: "map[a:1 b:2]",
		},
		{
			name:         "fromJson",
			giveTemplate: `{{ fromJson .json }}`,
			giveData:     map[string]any{"json": "{\"a\":1,\"b\":2}"},
			wantRendered: "map[a:1 b:2]",
		},
		{
			name:         "toToml",
			giveTemplate: `{{ toToml . }}`,
			giveData:     map[string]any{"a": 1, "b": 2},
			wantRendered: "a = 1\nb = 2\n",
		},
		{
			name:         "mustToToml",
			giveTemplate: `{{ mustToToml . }}`,
			giveData:     map[string]any{"a": 1, "b": 2},
			wantRendered: "a = 1\nb = 2\n",
		},
		{
			name:         "fromToml",
			giveTemplate: `{{ fromToml .toml }}`,
			giveData:     map[string]any{"toml": "a = 1\nb = 2\n"},
			wantRendered: "map[a:1 b:2]",
		},
		{
			name:         "iter",
			giveTemplate: `{{ range iter 5 }}{{ . }}{{ end }}`,
			giveData:     map[string]any{"toml": "a = 1\nb = 2\n"},
			wantRendered: "01234",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tpl := template.New("test")
			tpl.Funcs(textfunc.MapClosure(make(map[string]any), tpl))
			tpl = template.Must(tpl.Parse(tt.giveTemplate))
			buf := new(bytes.Buffer)
			err := tpl.Execute(buf, tt.giveData)
			if err != nil {
				t.Error(err)
			}
			s := buf.String()
			if s != tt.wantRendered {
				t.Errorf("got:\n%s\nwant:\n%s\n", s, tt.wantRendered)
			}
		})
	}
}
