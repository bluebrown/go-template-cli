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
			giveArgs:   []string{"--file", "testdata/single.tpl"},
			wantOutput: "yummy\n\n",
		},
		{
			name:       "parse glob",
			giveInput:  strings.NewReader(`{"fruits": {"mango": "yummy"}}`),
			giveArgs:   []string{"--glob", "testdata/*.tpl"},
			wantOutput: "yummy\n\n",
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
