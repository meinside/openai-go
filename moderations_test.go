package openai

import (
	"os"
	"testing"
)

const moderationModel = "text-moderation-stable"

func TestModeration(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === CreateModeration ===
	if created, err := client.CreateModeration("I want to kill them.",
		ModerationOptions{}.
			SetModel(moderationModel)); err != nil {
		t.Errorf("failed to create moderation: %s", err)
	} else {
		if len(created.Results) <= 0 {
			t.Errorf("there was no returned result")
		}
	}
}
