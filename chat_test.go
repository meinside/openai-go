package openai

import (
	"log"
	"os"
	"testing"
)

// const chatCompletionModel = "gpt-3.5-turbo"
const chatCompletionModel = "gpt-3.5-turbo-0613" // for testing `function`

// === CreateChatCompletion ===
func TestChatCompletions(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	if created, err := client.CreateChatCompletion(chatCompletionModel,
		[]ChatMessage{NewChatUserMessage("Hello!")},
		nil); err != nil {
		t.Errorf("failed to create chat completion: %s", err)
	} else {
		if len(created.Choices) <= 0 {
			t.Errorf("there was no returned choice")
		}
	}
}

// === CreateChatCompletion (stream) ===
func TestChatCompletionsStream(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

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
						log.Printf("stream response = %s", *comp.response.Choices[0].Delta.Content)
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

// === CreateChatCompletion (function) ===
func TestChatCompletionsFunction(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	if created, err := client.CreateChatCompletion(chatCompletionModel,
		[]ChatMessage{NewChatUserMessage("What's the weather like in Seoul?")},
		ChatCompletionOptions{}.
			SetFunctions([]ChatCompletionFunction{
				NewChatCompletionFunction(
					"get_current_weather",
					"Get the current weather in a given location",
					NewChatCompletionFunctionParameters().
						AddPropertyWithDescription("location", "string", "The city and state, e.g. San Francisco, CA").
						AddPropertyWithEnums("unit", "string", []string{"celsius", "fahrenheit"}).
						SetRequiredParameters([]string{"location", "unit"}),
				),
			}).
			SetFunctionCall(ChatCompletionFunctionCallAuto)); err != nil {
		t.Errorf("failed to create chat completion: %s", err)
	} else {
		if len(created.Choices) <= 0 {
			t.Errorf("there was no returned choice")
		} else {
			message := created.Choices[0].Message

			if message.FunctionCall == nil {
				t.Errorf("there was no returned function call")
			} else {
				functionName := message.FunctionCall.Name
				if functionName == "" {
					t.Errorf("there was no returned function call name")
				}

				if message.FunctionCall.Arguments == nil {
					t.Errorf("there were no returned function call arguments")
				} else {
					arguments, _ := message.FunctionCall.ArgumentsParsed()

					var location, unit string
					if l, exists := arguments["location"]; exists {
						location = l.(string)
					} else {
						t.Errorf("there was no returned parameter 'location' from function call")
					}
					if u, exists := arguments["unit"]; exists {
						unit = u.(string)
					} else {
						t.Errorf("there was no returned parameter 'unit' from function call")
					}

					t.Logf("will call %s(\"%s\", \"%s\")", functionName, location, unit)
					//functionResponse := `functionName`(location, unit) // -> get_current_weather('Seoul', 'celsius')
					functionResponse := "36.5"

					// FIXME: workaround for error: `'content' is a required property - 'messages.1'`
					content := "test"
					message.Content = &content

					if created, err := client.CreateChatCompletion(chatCompletionModel, []ChatMessage{
						NewChatUserMessage("What's the weather like in Seoul?"),
						message,
						NewChatFunctionMessage(functionName, functionResponse),
					}, nil); err != nil {
						t.Errorf("failed to create chat completion with local function response: %s", err)
					} else {
						if len(created.Choices) <= 0 {
							t.Errorf("there was no returned choice for local function response")
						}
					}
				}
			}
		}
	}
}
