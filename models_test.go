package openai

import (
	"os"
	"testing"
)

func TestModels(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === ListModels ===
	if models, err := client.LitModels(); err != nil {
		t.Errorf("failed to list models: %s", err)
	} else {
		if len(models.Data) > 0 {
			id := models.Data[0].ID
			// === RetrieveModel ===
			if _, err := client.RetrieveModel(id); err != nil {
				t.Errorf("failed to retrieve model: %s", err)
			}
		} else {
			t.Errorf("there was no returned item")
		}
	}
}
