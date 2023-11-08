package openai

import (
	"log"
	"os"
	"testing"
)

const chatCompletionModel = "gpt-3.5-turbo"
const chatCompletionVisionModel = "gpt-4-vision-preview"

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
						if content, err := comp.response.Choices[0].Delta.ContentString(); err == nil {
							log.Printf("stream response string: %s", content)
							continue
						}

						if content, err := comp.response.Choices[0].Delta.ContentArray(); err == nil {
							log.Printf("stream response messages: %+v", content)
							continue
						}
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
//
// example from: https://platform.openai.com/docs/guides/function-calling/parallel-function-calling
func TestChatCompletionsFunction(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	messages := []ChatMessage{
		NewChatUserMessage("What's the weather like in Seoul?"),
	}

	// generate a chat completion with function calls
	if created, err := client.CreateChatCompletion(chatCompletionModel,
		messages,
		ChatCompletionOptions{}.
			SetTools([]ChatCompletionTool{
				NewChatCompletionTool(
					"get_current_weather",
					"Get the current weather in a given location",
					NewChatCompletionFunctionParameters().
						AddPropertyWithDescription("location", "string", "The city and state, e.g. San Francisco, CA").
						AddPropertyWithEnums("unit", "string", []string{"celsius", "fahrenheit"}).
						SetRequiredParameters([]string{"location", "unit"}),
				),
			}).
			SetToolChoice(ChatCompletionToolChoiceAuto)); err != nil {
		t.Errorf("failed to create chat completion: %s", err)
	} else {
		if len(created.Choices) <= 0 {
			t.Errorf("there was no returned choice")
		} else {
			responseMessage := created.Choices[0].Message

			// FIXME: workaround for error: `'content' is a required property - 'messages.1'` <= assistant message's content should not be nil?
			content := "test"
			responseMessage.Content = &content

			// append the first response to the `messages`
			messages = append(messages, responseMessage)

			for _, toolCall := range responseMessage.ToolCalls {
				function := toolCall.Function

				// parse returned arguments into a struct
				type parsed struct {
					Location string `json:"location"`
					Unit     string `json:"unit"`
				}
				var arguments parsed
				if err := toolCall.ArgumentsInto(&arguments); err != nil {
					t.Errorf("failed to parse arguments into struct: %s", err)
				} else {
					t.Logf("will call %s(\"%s\", \"%s\")", function.Name, arguments.Location, arguments.Unit)

					// NOTE: get your local function's result with the generated arguments
					functionResponse := "36.5" //functionResponse := get_current_weather('Seoul', 'celsius')

					// and append it to the `messages`
					messages = append(messages, NewChatToolMessage(toolCall.ID, functionResponse))
				}
			}

			// generate a chat completion again with a local function result from the generated arguments
			if created, err := client.CreateChatCompletion(chatCompletionModel, messages, nil); err != nil {
				t.Errorf("failed to create chat completion with local function response: %s", err)
			} else {
				if len(created.Choices) <= 0 {
					t.Errorf("there was no returned choice for local function response")
				}
			}
		}
	}
}

// === CreateChatCompletion (function - stream) ===
//
// example from: https://platform.openai.com/docs/guides/function-calling/parallel-function-calling
func TestChatCompletionsFunctionStream(t *testing.T) {
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

	messages := []ChatMessage{
		NewChatUserMessage("What's the weather like in Seoul?"),
	}

	// generate a chat completion with function calls
	if _, err := client.CreateChatCompletion(chatCompletionModel,
		messages,
		ChatCompletionOptions{}.
			SetTools([]ChatCompletionTool{
				NewChatCompletionTool(
					"get_current_weather",
					"Get the current weather in a given location",
					NewChatCompletionFunctionParameters().
						AddPropertyWithDescription("location", "string", "The city and state, e.g. San Francisco, CA").
						AddPropertyWithEnums("unit", "string", []string{"celsius", "fahrenheit"}).
						SetRequiredParameters([]string{"location", "unit"}),
				),
			}).
			SetToolChoice(ChatCompletionToolChoiceAuto).
			SetStream(func(response ChatCompletion, done bool, err error) {
				ch <- completion{response: response, done: done, err: err}
				if done {
					toolCall := response.Choices[0].Message.ToolCalls[0]
					function := toolCall.Function

					// parse returned arguments into a struct
					type parsed struct {
						Location string `json:"location"`
						Unit     string `json:"unit"`
					}
					var arguments parsed
					if err := toolCall.ArgumentsInto(&arguments); err != nil {
						t.Errorf("failed to parse arguments into struct: %s", err)
					} else {
						t.Logf("will call %s(\"%s\", \"%s\")", function.Name, arguments.Location, arguments.Unit)

						// NOTE: get your local function's result with the generated arguments
					}

					close(ch)
				}
			})); err != nil {
		t.Errorf("failed to create chat completion with stream: %s", err)
	} else {
		for comp := range ch {
			if comp.err == nil {
				if client.Verbose {
					if !comp.done {
						if len(comp.response.Choices[0].Delta.ToolCalls) > 0 {
							toolCall := comp.response.Choices[0].Delta.ToolCalls[0]
							function := toolCall.Function

							if function.Name != "" {
								log.Printf("stream response function name: %s", function.Name)
								continue
							}

							if function.Arguments != "" {
								log.Printf("stream response function arguments (partial): %s", function.Arguments)
								continue
							}
						}
					}
				}
			} else {
				t.Errorf("there was an error in response stream: %s", comp.err)
			}
		}
	}
}

// === CreateChatCompletion (vision) ===
//
// https://platform.openai.com/docs/guides/vision/vision
func TestChatCompletionsVision(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	if image, err := NewFileParamFromFilepath("./sample/pepe.png"); err == nil {
		if created, err := client.CreateChatCompletion(chatCompletionVisionModel,
			[]ChatMessage{
				NewChatUserMessage[[]ChatMessageContent]([]ChatMessageContent{
					NewChatMessageContentWithText("What’s in this image?"),
					//NewChatMessageContentWithImageURL("https://user-images.githubusercontent.com/185988/60949207-a9daf400-a32f-11e9-8f11-d68d31cb0c31.png"),
					NewChatMessageContentWithFileParam(image),
				}),
			},
			nil); err != nil {
			t.Errorf("failed to create chat completion for vision: %s", err)
		} else {
			if len(created.Choices) <= 0 {
				t.Errorf("there was no returned choice")
			}
		}
	} else {
		t.Errorf("failed to open sample image: %s", err)
	}
}
