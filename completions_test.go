package openai

import (
	"os"
	"testing"
)

const (
	completionModel = "text-davinci-003"
)

func TestCompletions(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === CreateCompletion ===
	if created, err := client.CreateCompletion(completionModel,
		CompletionOptions{}.
			SetPrompt("Say this is a test.").
			SetMaxTokens(7).
			SetTemperature(0).
			SetTopP(1).
			SetN(1).
			SetStop("\n")); err != nil {
		t.Errorf("failed to create completion: %s", err)
	} else {
		if len(created.Choices) <= 0 {
			t.Errorf("there was no returned choice")
		}
	}
}
