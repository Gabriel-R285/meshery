
package model

import (
	"flag"
	"reflect"
	"testing"

	"github.com/meshery/meshery/server/models"
	display "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesherymodel "github.com/meshery/schemas/models/v1beta1/model"
)

var update = flag.Bool("update", false, "Update the model")

func TestGenerateModelOutput(t *testing.T) {
	modelsResponse := &models.MeshmodelsAPIResponse{
		Models: []mesherymodel.ModelDefinition{
			{
				Name:    "Prometheus",
				Version: "1.0.0",
			},
			{
				Name:    "Istio",
				Version: "1.20.0",
			},
		},
	}

	got := generateModelOutput(modelsResponse)

	want := []display.ModelOutput{
		{
			Model:   "Prometheus",
			Version: "1.0.0",
		},
		{
			Model:   "Istio",
			Version: "1.20.0",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("generateModelOutput() = %v, want %v", got, want)
	}
}

func TestListModelOutputFlag(t *testing.T) {
	flag := listModelCmd.Flags().Lookup("output")

	if flag == nil {
		t.Fatal("expected output flag to be registered")
	}

	if flag.DefValue != "table" {
		t.Errorf("expected default output to be table, got %q", flag.DefValue)
	}
}
