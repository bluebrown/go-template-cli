package main

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"
)

func Test_commandLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		giveInput    string
		giveArgs     []string
		wantOutput   string
		wantErrRegex string
	}{
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
			name:         "decoder invalid",
			giveInput:    `foo = "bar"`,
			giveArgs:     []string{"--no-newline", "--decoder=yikes", `test: {{.foo}}`},
			wantErrRegex: "invalid argument",
		},
		{
			name:         "flag help",
			giveArgs:     []string{"-h"},
			wantErrRegex: "parse flags: pflag: help requested",
		},
		{
			name:       "parse file",
			giveInput:  `{"fruits": {"mango": "yummy"}}`,
			giveArgs:   []string{"--file", "testdata/mango.tpl"},
			wantOutput: "yummy\n\n",
		},
		{
			name:       "parse glob",
			giveInput:  `{"fruits": {"mango": "yummy", "apple": "yuk"}}`,
			giveArgs:   []string{"--glob", "testdata/*.tpl", "--name", "apple.tpl"},
			wantOutput: "yuk\n\n",
		},
		{
			name:         "require name",
			giveInput:    `{}`,
			giveArgs:     []string{"--glob", "testdata/*.tpl"},
			wantErrRegex: `the --name flag is required when multiple templates are defined`,
		},
		{
			name:         "require template",
			giveInput:    `{}`,
			wantErrRegex: "no templates found",
		},
		{
			name:       "no value default",
			giveInput:  `{}`,
			giveArgs:   []string{"{{.nope}}"},
			wantOutput: "<no value>\n",
		},
		{
			name:         "no value error",
			giveInput:    `{}`,
			giveArgs:     []string{"{{.nope}}", "--option", "missingkey=error"},
			wantErrRegex: `map has no entry for key "nope"`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			err := commandLine(context.Background(), tt.giveArgs, strings.NewReader(tt.giveInput), output)
			if len(tt.wantErrRegex) > 0 {
				if err == nil {
					t.Fatalf("want error but got none")
				}
				ok, err := regexp.MatchString(tt.wantErrRegex, err.Error())
				if err != nil {
					t.Fatal(err)
				}
				if !ok {
					t.Fatalf("wrong error: got %q but want %q", err.Error(), tt.wantErrRegex)
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
