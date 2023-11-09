package openai

import (
	"os"
	"testing"
)

const (
	speechModel = "tts-1"
	audioModel  = "whisper-1"
)

func TestAudio(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === CreateSpeech ===
	if speech, err := client.CreateSpeech(speechModel, "All your base are belong to us.", SpeechVoiceAlloy, nil); err != nil {
		t.Errorf("failed to create speech: %s", err)
	} else {
		if len(speech) <= 0 {
			t.Errorf("returned speech bytes is empty")
		}
	}

	// === CreateTranscription ===
	if audio, err := NewFileParamFromFilepath("./sample/test.mp3"); err != nil {
		t.Errorf("failed to open sample audio: %s", err)
	} else {
		if translated, err := client.CreateTranscription(audio, audioModel, nil); err != nil {
			t.Errorf("failed to create transcription: %s", err)
		} else {
			if translated.JSON == nil &&
				translated.Text == nil &&
				translated.SRT == nil &&
				translated.VerboseJSON == nil &&
				translated.VTT == nil {
				t.Errorf("there was no returned data")
			}
		}
	}

	// === CreateTranslation ===
	if audio, err := NewFileParamFromFilepath("./sample/test.mp3"); err != nil {
		t.Errorf("failed to open sample audio: %s", err)
	} else {
		if translated, err := client.CreateTranslation(audio, audioModel, nil); err != nil {
			t.Errorf("failed to create translation: %s", err)
		} else {
			if translated.JSON == nil &&
				translated.Text == nil &&
				translated.SRT == nil &&
				translated.VerboseJSON == nil &&
				translated.VTT == nil {
				t.Errorf("there was no returned data")
			}
		}
	}
}
