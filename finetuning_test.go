package openai

import (
	"os"
	"testing"
)

func TestFineTuning(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	if file, err := NewFileParamFromFilepath("./sample/training.jsonl"); err != nil {
		t.Errorf("failed to open sample jsonl file: %s", err)
	} else {
		if uploaded, err := client.UploadFile(file, "fine-tune"); err != nil {
			t.Errorf("failed to upload file: %s", err)
		} else {
			// === CreateFineTuningJob ===
			if created, err := client.CreateFineTuningJob(uploaded.ID, "davinci-002", nil); err != nil {
				t.Errorf("failed to create fine-tuning job: %s", err)
			} else {
				// === ListFineTuningJobs ===
				if _, err := client.ListFineTuningJobs(nil); err != nil {
					t.Errorf("failed to list fine-tuning jobs: %s", err)
				}

				// === RetrieveFineTuningJob ===
				if retrieved, err := client.RetrieveFineTuningJob(created.ID); err != nil {
					t.Errorf("failed to retrieve a fine-tuning job: %s", err)
				} else {
					// === ListFineTuningJobEvents ===
					if events, err := client.ListFineTuningJobEvents(retrieved.ID, nil); err != nil {
						t.Errorf("failed to list fine-tuning job events: %s", err)
					} else {
						numFineTuningJobEvents := len(events.Data)

						if numFineTuningJobEvents <= 0 {
							t.Errorf("there was no fine-tuning job events")
						} else {
							// === CancelFineTuningJob ===
							if cancelled, err := client.CancelFineTuningJob(retrieved.ID); err != nil {
								t.Errorf("failed to cancel fine-tuning job: %s", err)
							} else {
								if cancelled.ID != retrieved.ID {
									t.Errorf("canceled fine-tuning job id is different from the requested one: '%s' vs '%s'", cancelled.ID, retrieved.ID)
								}
							}
						}
					}
				}
			}
		}
	}
}
