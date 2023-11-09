package openai

import (
	"os"
	"testing"
)

func TestThreads(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	client.SetBetaHeader(`assistants=v1`)

	// === CreateThread ===
	if created, err := client.CreateThread(CreateThreadOptions{}.
		SetMessages([]ThreadMessage{
			NewThreadMessage("What's the weather like in Seoul, Korea?"),
		})); err != nil {
		t.Errorf("failed to create thread: %s", err)
	} else {
		threadID := created.ID

		// === RetrieveThread ===
		if retrieved, err := client.RetrieveThread(threadID); err != nil {
			t.Errorf("failed to retrieve thread: %s", err)
		} else {
			if retrieved.ID != threadID {
				t.Errorf("retrieved thread's id is not equal to the created one's")
			}

			// === ModifyThread ===
			if modified, err := client.ModifyThread(threadID, ModifyThreadOptions{}.SetMetadata(map[string]string{
				"country": "Korea",
				"city":    "Seoul",
			})); err != nil {
				t.Errorf("failed to modify thread: %s", err)
			} else {
				if len(modified.Metadata) != 2 {
					t.Errorf("modified metadata count is not 2")
				}
			}

			// === DeleteThread ===
			if deleted, err := client.DeleteThread(threadID); err != nil {
				t.Errorf("failed to delete thread: %s", err)
			} else {
				if !deleted.Deleted {
					t.Errorf("deleted status of deleted thread is not true")
				}
			}
		}
	}
}
