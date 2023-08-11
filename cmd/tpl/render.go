package main

import (
	"errors"
	"fmt"
	"io"
)

// render a template
func (cli *state) render(w io.Writer, data any) error {
	templateName, err := cli.selectTemplate()
	if err != nil {
		return fmt.Errorf("select template: %w", err)
	}

	if err := cli.template.ExecuteTemplate(w, templateName, data); err != nil {
		return fmt.Errorf("execute template: %v", err)
	}

	if !cli.noNewline {
		fmt.Fprintln(w)
	}

	return nil
}

// determine the template to execute. In the order of precedence:
//  1. current name, if set
//  2. default name, if at least 1 positional arg
//  3. templates name, if exactly 1 template
//  4. --name flag required, otherwise
func (cli *state) selectTemplate() (string, error) {
	templates := cli.template.Templates()

	if len(templates) == 0 {
		return "", errors.New("no templates found")
	}

	if cli.templateName != "" {
		return cli.templateName, nil
	}

	if len(cli.posArgs) > 0 {
		return cli.defaultTemplateName, nil
	}

	if len(templates) == 1 {
		return templates[0].Name(), nil
	}

	return "", fmt.Errorf("the --name flag is required when multiple templates are defined and no default template exists")
}
