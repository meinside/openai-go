package openai

import (
	"log"
	"os"
	"testing"
)

const chatCompletionModel = "gpt-3.5-turbo"

func TestChatCompletions(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === CreateChatCompletion ===
	if created, err := client.CreateChatCompletion(chatCompletionModel,
		[]ChatMessage{NewChatUserMessage("Hello!")},
		nil); err != nil {
		t.Errorf("failed to create chat completion: %s", err)
	} else {
		if len(created.Choices) <= 0 {
			t.Errorf("there was no returned choice")
		}
	}

	// === CreateChatCompletion (stream) ===
	type completion struct {
		response ChatCompletion
		done     bool
		err      error
	}
	ch := make(chan completion, 1)
	if _, err := client.CreateChatCompletion(chatCompletionModel,
		[]ChatMessage{NewChatUserMessage("Hello!")},
		ChatCompletionOptions{}.
			SetStream(func(response ChatCompletion, done bool, err error) {
				ch <- completion{response: response, done: done, err: err}
				if done {
					close(ch)
				}
			})); err == nil {
		for comp := range ch {
			if comp.err == nil {
				if client.Verbose {
					if !comp.done {
						log.Printf("stream response = %s", comp.response.Choices[0].Delta.Content)
					}
				}
			} else {
				t.Errorf("there was an error in response stream: %s", comp.err)
			}
		}
	} else {
		t.Errorf("failed to create chat completion with stream: %s", err)
	}
}
