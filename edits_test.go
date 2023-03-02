package openai

import (
	"os"
	"testing"
)

const editModel = "text-davinci-edit-001"

func TestEdit(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === CreateEdit ===
	if created, err := client.CreateEdit(editModel,
		"Fix the spelling mistakes",
		EditOptions{}.
			SetInput("What day of the wek is it?")); err != nil {
		t.Errorf("failed to create edit: %s", err)
	} else {
		if len(created.Choices) <= 0 {
			t.Errorf("there was no returned choice")
		}
	}
}
