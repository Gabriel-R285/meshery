package display

import (
	"bytes"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	var output bytes.Buffer

	p := Printer{
		Writer: &output,
		Format: OutputJSON,
	}

	models := []ModelOutput{
		{
			Model:   "kubernetes",
			Version: "v1.30",
		},
	}

	err := p.Print(models)
	if err != nil {
		t.Fatalf("Print() returned error: %v", err)
	}

	expected := `[
  {
    "model": "kubernetes",
    "version": "v1.30"
  }
]
`

	if output.String() != expected {
		t.Errorf(
			"unexpected output\nexpected:\n%s\ngot:\n%s",
			expected,
			output.String(),
		)
	}
}

func TestPrintYAML(t *testing.T) {
	var output bytes.Buffer

	p := Printer{
		Writer: &output,
		Format: OutputYAML,
	}

	models := []ModelOutput{
		{
			Model:   "kubernetes",
			Version: "v1.30",
		},
	}

	err := p.Print(models)
	if err != nil {
		t.Fatalf("Print() returned error: %v", err)
	}

	expected := `- model: kubernetes
  version: v1.30
`

	if output.String() != expected {
		t.Errorf(
			"unexpected output\nexpected:\n%s\ngot:\n%s",
			expected,
			output.String(),
		)
	}
}

func TestUnsupportedFormat(t *testing.T) {
	var output bytes.Buffer

	p := Printer{
		Writer: &output,
		Format: OutputFormat("xml"),
	}

	err := p.Print([]ModelOutput{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}