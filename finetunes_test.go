package openai

import (
	"os"
	"testing"
)

func TestFineTune(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	if file, err := NewFileParamFromFilepath("./sample/training.jsonl"); err == nil {
		if uploaded, err := client.UploadFile(file, "fine-tune"); err != nil {
			t.Errorf("failed to upload file: %s", err)
		} else {
			// === CreateFineTune ===
			if created, err := client.CreateFineTune(uploaded.ID, nil); err != nil {
				t.Errorf("failed to create fine-tune: %s", err)
			} else {
				if len(created.Events) <= 0 {
					t.Errorf("there was no returned event")
				} else {
					// === ListFineTunes ===
					if fineTunes, err := client.ListFineTunes(); err != nil {
						t.Errorf("failed to list fine-tunes: %s", err)
					} else {
						// === RetrieveFineTune ===
						numFineTunes := len(fineTunes.Data)
						if numFineTunes <= 0 {
							t.Errorf("there was no fine-tune")
						} else {
							lastFineTuneID := fineTunes.Data[numFineTunes-1].ID
							if fineTune, err := client.RetrieveFineTune(lastFineTuneID); err != nil {
								t.Errorf("failed to retrieve fine-tune: %s", err)
							} else {
								// === ListFineTuneEvents ===
								if events, err := client.ListFineTuneEvents(fineTune.ID, nil); err != nil {
									t.Errorf("failed to list fine-tune events: %s", err)
								} else {
									if len(events.Data) <= 0 {
										t.Errorf("there was no returned event")
									} else {
										// === CancelFineTune ===
										if canceled, err := client.CancelFineTune(fineTune.ID); err != nil {
											t.Errorf("failed to cancel fine-tune: %s", err)
										} else {
											if canceled.ID != fineTune.ID {
												t.Errorf("canceled fine-tune's id does not match the requested one: %s - %s", canceled.ID, fineTune.ID)
											}
										}

										// === DeleteFineTuneModel ===
										// TODO
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
