package openai

import (
	"os"
	"testing"
)

const embeddingModel = "text-embedding-3-small"

func TestEmbeddings(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === CreateEmbedding ===
	if created, err := client.CreateEmbedding(embeddingModel, "The food was delicious and the waiter...", nil); err != nil {
		t.Errorf("failed to create embedding: %s", err)
	} else {
		if len(created.Data) <= 0 {
			t.Errorf("there were no returned item")
		}
	}
}
