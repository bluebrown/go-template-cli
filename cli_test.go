package cli

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func Test_state_run(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		giveInput      string
		giveArgs       []string
		wantOutput     string
		wantErrorMatch string
	}{
		{
			name:       "parse arg",
			giveInput:  `{"fruits": {"mango": "yummy"}}`,
			giveArgs:   []string{"{{.fruits.mango}}"},
			wantOutput: "yummy\n",
		},
		{
			name:       "parse file",
			giveInput:  `{"fruits": {"mango": "yummy"}}`,
			giveArgs:   []string{"--file", "testdata/mango.tpl"},
			wantOutput: "yummy\n\n",
		},
		{
			name:       "parse glob single",
			giveInput:  `{"fruits": {"mango": "yummy"}}`,
			giveArgs:   []string{"--glob", "testdata/ma*.tpl"},
			wantOutput: "yummy\n\n",
		},
		{
			name:       "parse multi with positional",
			giveInput:  `{"fruits": {"kiwi": "yay"}}`,
			giveArgs:   []string{"{{.fruits.kiwi}}", "--file", "testdata/apple.tpl", "--glob", "testdata/ma*.tpl"},
			wantOutput: "yay\n",
		},
		{
			name:       "parse multi without positional",
			giveInput:  `{"fruits": {"avocado": "impossibru"}}`,
			giveArgs:   []string{"--glob", "testdata/*.tpl", "--name", "avocado.tpl"},
			wantOutput: "impossibru\n\n",
		},
		{
			name:           "require name",
			giveInput:      `{}`,
			giveArgs:       []string{"--glob", "testdata/*.tpl"},
			wantErrorMatch: `the --name flag is required when multiple templates are defined`,
		},
		{
			name:           "require template",
			giveInput:      `{}`,
			wantErrorMatch: "no templates found",
		},
		{
			name:           "flag help",
			giveArgs:       []string{"-h"},
			wantErrorMatch: "pflag: help requested",
		},
		{
			name:       "no value default",
			giveInput:  `{}`,
			giveArgs:   []string{"{{.nope}}"},
			wantOutput: "<no value>\n",
		},
		{
			name:           "no value error",
			giveInput:      `{}`,
			giveArgs:       []string{"{{.nope}}", "--option", "missingkey=error"},
			wantErrorMatch: `map has no entry for key "nope"`,
		},
		{
			name:       "decoder yaml",
			giveInput:  "foo: bar",
			giveArgs:   []string{"--no-newline", "--decoder=yaml", `test: {{.foo}}`},
			wantOutput: "test: bar",
		},
		{
			name:       "decoder json",
			giveInput:  `{"foo":"bar"}`,
			giveArgs:   []string{"--no-newline", "--decoder=json", `test: {{.foo}}`},
			wantOutput: "test: bar",
		},
		{
			name:       "decoder toml",
			giveInput:  `foo = "bar"`,
			giveArgs:   []string{"--no-newline", "--decoder=toml", `test: {{.foo}}`},
			wantOutput: "test: bar",
		},
		{
			name:           "decoder invalid",
			giveInput:      `foo = "bar"`,
			giveArgs:       []string{"--no-newline", "--decoder=yikes", `test: {{.foo}}`},
			wantErrorMatch: `unsupported decoder "yikes"`,
		},
		{
			name:       "parse file with equal flag",
			giveInput:  `{"fruits": {"mango": "yummy"}}`,
			giveArgs:   []string{"--file=testdata/mango.tpl"},
			wantOutput: "yummy\n\n",
		},
		{
			name:       "yaml string keys",
			giveInput:  `key: value`,
			giveArgs:   []string{"--no-newline", "-d=yaml", `{{printf "%#v" .}}`},
			wantOutput: `map[string]interface {}{"key":"value"}`,
		},
		{
			name:       "nil data",
			giveInput:  "",
			giveArgs:   []string{`{{.}}`},
			wantOutput: "<no value>\n",
		},
		{
			name:       "issue-19-1",
			giveInput:  `{"fruits": {"mango": "yummy"}}`,
			giveArgs:   []string{"-f", "testdata/apple.tpl", "-f", "testdata/mango.tpl", "-n", "mango.tpl"},
			wantOutput: "yummy\n\n",
		},
		{
			name:           "issue-19-2",
			giveArgs:       []string{"-f", "testdata/apple.tpl", "-f", "testdata/mango.tpl"},
			wantErrorMatch: `the --name flag is required when multiple templates are defined`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := &bytes.Buffer{}

			var in io.Reader
			if tt.giveInput != "" {
				in = strings.NewReader(tt.giveInput)
			}

			err := New("test", pflag.NewFlagSet(tt.name, pflag.ContinueOnError)).Run(tt.giveArgs, in, output)

			if len(tt.wantErrorMatch) > 0 {
				if err == nil {
					t.Fatalf("want error but got none")
				}
				ok, err2 := regexp.MatchString(tt.wantErrorMatch, err.Error())
				if err2 != nil {
					t.Fatal(err2)
				}
				if !ok {
					t.Fatalf("wrong error: got %v but want %q", err, tt.wantErrorMatch)
				}
				return
			} else if err != nil {
				t.Fatal(err)
			} else {
				if gotOutput := output.String(); gotOutput != tt.wantOutput {
					t.Errorf("wrong run output:\ngot:\n%q\n\nwant\n%q\n", gotOutput, tt.wantOutput)
				}
			}
		})
	}
}
