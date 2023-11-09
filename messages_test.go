package openai

import (
	"os"
	"testing"
)

func TestMessages(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	client.SetBetaHeader(`assistants=v1`)

	if thread, err := client.CreateThread(nil); err != nil {
		t.Errorf("failed to create a thread for testing messages: %s", err)
	} else {
		threadID := thread.ID

		if file, err := NewFileParamFromFilepath("./sample/test.rb"); err != nil {
			t.Errorf("failed to open file for testing thread message: %s", err)
		} else {
			if uploaded, err := client.UploadFile(file, "assistants"); err != nil {
				t.Errorf("failed to upload file for testing thread message: %s", err)
			} else {
				// === CreateMessage ===
				if created, err := client.CreateMessage(threadID, "user", "What is the weather like in Seoul, Korea?", CreateMessageOptions{}.
					SetFileIDs([]string{
						uploaded.ID,
					})); err != nil {
					t.Errorf("failed to create thread message: %s", err)
				} else {
					messageID := created.ID

					// === ListMessages ===
					if listed, err := client.ListMessages(threadID, ListMessagesOptions{}); err != nil {
						t.Errorf("failed to list thread messages: %s", err)
					} else {
						if len(listed.Data) <= 0 {
							t.Errorf("there was no returned thread message")
						}
					}

					// === RetrieveMessage ===
					if retrieved, err := client.RetrieveMessage(threadID, messageID); err != nil {
						t.Errorf("failed to retrieve thread message: %s", err)
					} else {
						if messageID != retrieved.ID {
							t.Errorf("retrieved message id: %s does not match the requested one: %s", retrieved.ID, messageID)
						}

						// === ModifyMessage ===
						if modified, err := client.ModifyMessage(threadID, messageID, ModifyMessageOptions{}.SetMetadata(map[string]string{})); err != nil {
							t.Errorf("failed to modify thread message: %s", err)
						} else {
							if modified.ID != messageID {
								t.Errorf("modified message id: %s does not match the requsted one: %s", modified.ID, messageID)
							}

							// === ListMessageFiles ===
							if files, err := client.ListMessageFiles(threadID, messageID, nil); err != nil {
								t.Errorf("failed to list message files: %s", err)
							} else {
								file := files.Data[0]

								// === RetrieveMessageFile ===
								if _, err := client.RetrieveMessageFile(threadID, messageID, file.ID); err != nil {
									t.Errorf("failed to retrieve message file: %s", err)
								}
							}
						}
					}
				}
			}
		}
	}
}
