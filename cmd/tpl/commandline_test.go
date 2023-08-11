package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func Test_commandLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		giveInput  io.Reader
		giveArgs   []string
		wantOutput string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "decoder yaml",
			giveInput:  strings.NewReader("foo: bar"),
			giveArgs:   []string{"--no-newline", "--decoder=yaml", `test: {{.foo}}`},
			wantOutput: "test: bar",
		},
		{
			name:       "decoder json",
			giveInput:  strings.NewReader(`{"foo":"bar"}`),
			giveArgs:   []string{"--no-newline", "--decoder=json", `test: {{.foo}}`},
			wantOutput: "test: bar",
		},
		{
			name:       "decoder toml",
			giveInput:  strings.NewReader(`foo = "bar"`),
			giveArgs:   []string{"--no-newline", "--decoder=toml", `test: {{.foo}}`},
			wantOutput: "test: bar",
		},
		{
			name:       "decoder invalid",
			giveInput:  strings.NewReader(`foo = "bar"`),
			giveArgs:   []string{"--no-newline", "--decoder=yikes", `test: {{.foo}}`},
			wantErr:    true,
			wantErrMsg: "parse flags: invalid argument \"yikes\" for \"-d, --decoder\" flag: invalid decoder kind: yikes, supported value are: json, yaml, toml",
		},
		{
			name:       "flag help",
			giveArgs:   []string{"-h"},
			wantErr:    true,
			wantErrMsg: "parse flags: pflag: help requested",
		},
		{
			name:       "parse file",
			giveInput:  strings.NewReader(`{"fruits": {"mango": "yummy"}}`),
			giveArgs:   []string{"--file", "testdata/mango.tpl"},
			wantOutput: "yummy\n\n",
		},
		{
			name:       "parse glob",
			giveInput:  strings.NewReader(`{"fruits": {"mango": "yummy", "apple": "yuk"}}`),
			giveArgs:   []string{"--glob", "testdata/*.tpl", "--name", "apple.tpl"},
			wantOutput: "yuk\n\n",
		},
		{
			name:       "require name",
			giveInput:  strings.NewReader(`{}`),
			giveArgs:   []string{"--glob", "testdata/*.tpl"},
			wantErr:    true,
			wantErrMsg: `the --name flag is required when multiple templates are defined and no default template exists; defined templates are: "apple.tpl", "mango.tpl"`,
		},
		{
			name:       "require template",
			giveInput:  strings.NewReader(`{}`),
			wantErr:    true,
			wantErrMsg: "no templates found",
		},
		{
			name:       "no value default",
			giveInput:  strings.NewReader(`{}`),
			giveArgs:   []string{"{{.nope}}"},
			wantOutput: "<no value>\n",
		},
		{
			name:       "no value error",
			giveInput:  strings.NewReader(`{}`),
			giveArgs:   []string{"{{.nope}}", "--option", "missingkey=error"},
			wantErr:    true,
			wantErrMsg: `error executing template: template: _gotpl_default:1:2: executing "_gotpl_default" at <.nope>: map has no entry for key "nope"`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			err := commandLine(context.Background(), tt.giveArgs, tt.giveInput, output)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && (err.Error() != tt.wantErrMsg) {
				t.Errorf("run() error:\n%s\n\nwantErrMsg:\n%s\n", err.Error(), tt.wantErrMsg)
			}

			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("run() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
