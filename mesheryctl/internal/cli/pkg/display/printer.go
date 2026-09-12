package display

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type OutputFormat string

const (
	OutputJSON OutputFormat = "json"
	OutputYAML OutputFormat = "yaml"
)

type ModelOutput struct {
	Model   string `json:"model" yaml:"model"`
	Version string `json:"version" yaml:"version"`
}

type Printer struct {
	Writer io.Writer
	Format OutputFormat
}

func (p Printer) Print(models []ModelOutput) error {
	switch p.Format {
	case OutputJSON:
		return p.printJSON(models)

	case OutputYAML:
		return p.printYAML(models)

	default:
		return fmt.Errorf("unsupported output format: %s", p.Format)
	}
}

func (p Printer) printJSON(models []ModelOutput) error {
	data, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(p.Writer, string(data))
	return err
}

func (p Printer) printYAML(models []ModelOutput) error {
	data, err := yaml.Marshal(models)
	if err != nil {
		return err
	}

	_, err = p.Writer.Write(data)
	return err
}