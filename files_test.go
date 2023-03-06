package openai

import (
	"log"
	"os"
	"testing"
)

func TestFiles(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === ListFiles ===
	if _, err := client.ListFiles(); err != nil {
		t.Errorf("failed to list files: %s", err)
	}

	// === UploadFile ===
	if file, err := NewFileParamFromFilepath("./sample/training.jsonl"); err == nil {
		if uploaded, err := client.UploadFile(file, "fine-tune"); err != nil {
			t.Errorf("failed to upload file: %s", err)
		} else {
			fileID := uploaded.ID

			// === RetrieveFile ===
			if retrieved, err := client.RetrieveFile(fileID); err != nil {
				t.Errorf("failed to retrieve file: %s", err)
			} else {
				if retrieved.ID != fileID {
					t.Errorf("retrieved file's id does not match the requested one: %s - %s", retrieved.ID, fileID)
				}
			}

			// === RetrieveFileContent ===
			if bytes, err := client.RetrieveFileContent(fileID); err != nil {
				t.Errorf("failed to retrieve content of file: %s", err)
			} else {
				if len(bytes) != len(file.bs) {
					// test
					log.Printf("bytes = %s", string(bytes))

					t.Errorf("retrieved file content's bytes count does not match the original one: %d - %d", len(bytes), len(file.bs))
				}
			}

			// === DeleteFile ===
			if deleted, err := client.DeleteFile(fileID); err != nil {
				t.Errorf("failed to delete file: %s", err)
			} else {
				if deleted.ID != fileID {
					t.Errorf("deleted file's id does not match the requested one: %s - %s", deleted.ID, fileID)
				}
			}
		}
	} else {
		t.Errorf("failed to open sample jsonl file: %s", err)
	}
}
