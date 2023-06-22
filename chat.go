package openai

// https://platform.openai.com/docs/api-reference/chat

import (
	"encoding/json"
	"fmt"
)

// ChatMessageRole type for constants
type ChatMessageRole string

const (
	ChatMessageRoleSystem    ChatMessageRole = "system"
	ChatMessageRoleUser      ChatMessageRole = "user"
	ChatMessageRoleAssistant ChatMessageRole = "assistant"
	ChatMessageRoleFunction  ChatMessageRole = "function"
)

// ChatCompletionFunctionCall struct
type ChatCompletionFunctionCall struct {
	Name      string  `json:"name"`
	Arguments *string `json:"arguments,omitempty"` // = JSON string
}

// ArgumentsParsed returns the parsed map from ChatCompletionFunctionCall.
func (c ChatCompletionFunctionCall) ArgumentsParsed() (result map[string]any, err error) {
	if c.Arguments != nil {
		err = json.Unmarshal([]byte(*c.Arguments), &result)
	}

	return result, err
}

// ArgumentsInto parses the generated arguments into a given interface
func (c ChatCompletionFunctionCall) ArgumentsInto(out any) (err error) {
	if c.Arguments == nil {
		err = fmt.Errorf("parse failed: `arguments` is nil")
	} else {
		err = json.Unmarshal([]byte(*c.Arguments), &out)
	}

	return err
}

// ChatMessage struct for chat completion
//
// https://platform.openai.com/docs/guides/chat/introduction
type ChatMessage struct {
	Role    ChatMessageRole `json:"role"`
	Content *string         `json:"content,omitempty"`

	// for function call
	Name         *string                     `json:"name,omitempty"`
	FunctionCall *ChatCompletionFunctionCall `json:"function_call,omitempty"`
}

// ChatCompletionFunction struct for chat completion function
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-functions
type ChatCompletionFunction struct {
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// ChatCompletionFunctionParameters type
type ChatCompletionFunctionParameters map[string]any

// NewChatCompletionFunction returns a new chat completion function.
func NewChatCompletionFunction(name, description string, parameters ChatCompletionFunctionParameters) ChatCompletionFunction {
	return ChatCompletionFunction{
		Name:        name,
		Description: &description,
		Parameters:  parameters,
	}
}

// NewChatCompletionFunctionParameters returns an empty ChatCompletionFunctionParameters.
func NewChatCompletionFunctionParameters() ChatCompletionFunctionParameters {
	return ChatCompletionFunctionParameters{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
}

// AddPropertyWithDescription adds/overwrites a property in chat completion function parameters with a description.
func (p ChatCompletionFunctionParameters) AddPropertyWithDescription(name, typ3, description string) ChatCompletionFunctionParameters {
	if properties, exists := p["properties"]; exists {
		ps := properties.(map[string]any)
		ps[name] = map[string]string{
			"type":        typ3,
			"description": description,
		}
		p["properties"] = ps
	}
	return p
}

// AddPropertyWithEnums adds/overwrites a property in chat completion function parameters with enums.
func (p ChatCompletionFunctionParameters) AddPropertyWithEnums(name, typ3 string, enums []string) ChatCompletionFunctionParameters {
	if properties, exists := p["properties"]; exists {
		ps := properties.(map[string]any)
		ps[name] = map[string]any{
			"type": typ3,
			"enum": enums,
		}
		p["properties"] = ps
	}
	return p
}

// SetRequiredParameters sets/overwrites required parameter names for chat completion function parameters.
func (p ChatCompletionFunctionParameters) SetRequiredParameters(names []string) ChatCompletionFunctionParameters {
	p["required"] = names
	return p
}

// ChatCompletionFunctionCallMode type
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-function_call
type ChatCompletionFunctionCallMode string

// ChatCompletionFunctionCall constants
const (
	ChatCompletionFunctionCallNone ChatCompletionFunctionCallMode = "none"
	ChatCompletionFunctionCallAuto ChatCompletionFunctionCallMode = "auto"
)

// NewChatSystemMessage returns a new ChatMessage with system role.
func NewChatSystemMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleSystem,
		Content: &message,
	}
}

// NewChatUserMessage returns a new ChatMessage with user role.
func NewChatUserMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleUser,
		Content: &message,
	}
}

// NewChatAssistantMessage returns a new ChatMessage with assistant role.
func NewChatAssistantMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleAssistant,
		Content: &message,
	}
}

// NewChatFunctionMessage returns a new ChatMessage with function role.
func NewChatFunctionMessage(name, content string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleFunction,
		Name:    &name,
		Content: &content,
	}
}

// ChatCompletionChoice struct for chat completion response
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Delta        ChatMessage `json:"delta"` // Only appears in stream response
}

// ChatCompletion struct for chat completion response
type ChatCompletion struct {
	CommonResponse

	ID      string                 `json:"id"`
	Created int64                  `json:"created"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

// ChatCompletionOptions for creating chat completions
type ChatCompletionOptions map[string]any

// SetFunctions sets the `functions` parameter of chat completion request.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-functions
func (o ChatCompletionOptions) SetFunctions(functions []ChatCompletionFunction) ChatCompletionOptions {
	o["functions"] = functions
	return o
}

// SetFunctionCall sets the `function_call` parameter of chat completion request.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-function_call
func (o ChatCompletionOptions) SetFunctionCall(functionCall ChatCompletionFunctionCallMode) ChatCompletionOptions {
	o["function_call"] = functionCall
	return o
}

// SetFunctionCallWithName sets the `function_call` parameter of chat completion request.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-function_call
func (o ChatCompletionOptions) SetFunctionCallWithName(name string) ChatCompletionOptions {
	o["function_call"] = map[string]any{
		"name": name,
	}
	return o
}

// SetTemperature sets the `temperature` parameter of chat completion request.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-temperature
func (o ChatCompletionOptions) SetTemperature(temperature float64) ChatCompletionOptions {
	o["temperature"] = temperature
	return o
}

// SetTopP sets the `top_p` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-top_p
func (o ChatCompletionOptions) SetTopP(topP float64) ChatCompletionOptions {
	o["top_p"] = topP
	return o
}

// SetN sets the `n` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-n
func (o ChatCompletionOptions) SetN(n int) ChatCompletionOptions {
	o["n"] = n
	return o
}

// SetStream sets the `stream` parameter of chat completions.
//
// https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-stream
func (o ChatCompletionOptions) SetStream(cb callback) ChatCompletionOptions {
	o["stream"] = cb
	return o
}

// SetStop sets the `stop` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-stop
func (o ChatCompletionOptions) SetStop(stop any) ChatCompletionOptions {
	o["stop"] = stop
	return o
}

// SetMaxTokens sets the `max_tokens` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-max_tokens
func (o ChatCompletionOptions) SetMaxTokens(maxTokens int) ChatCompletionOptions {
	o["max_tokens"] = maxTokens
	return o
}

// SetPresencePenalty sets the `presence_penalty` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-presence_penalty
func (o ChatCompletionOptions) SetPresencePenalty(presencePenalty float64) ChatCompletionOptions {
	o["presence_penalty"] = presencePenalty
	return o
}

// SetFrequencyPenalty sets the `frequency_penalty` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-frequency_penalty
func (o ChatCompletionOptions) SetFrequencyPenalty(frequencyPenalty float64) ChatCompletionOptions {
	o["frequency_penalty"] = frequencyPenalty
	return o
}

// SetLogitBias sets the `logit_bias` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
func (o ChatCompletionOptions) SetLogitBias(logitBias map[string]any) ChatCompletionOptions {
	o["logit_bias"] = logitBias
	return o
}

// SetUser sets the `user` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-user
func (o ChatCompletionOptions) SetUser(user string) ChatCompletionOptions {
	o["user"] = user
	return o
}

// CreateChatCompletion creates a completion for chat messages.
//
// https://platform.openai.com/docs/api-reference/chat/create
func (c *Client) CreateChatCompletion(model string, messages []ChatMessage, options ChatCompletionOptions) (response ChatCompletion, err error) {
	if options == nil {
		options = ChatCompletionOptions{}
	}
	options["model"] = model
	options["messages"] = messages

	if options["stream"] != nil {
		cb := options["stream"].(callback)
		options["stream"] = true
		_, err := c.postCB("v1/chat/completions", options, cb)

		return ChatCompletion{}, err
	}

	var bytes []byte
	if bytes, err = c.post("v1/chat/completions", options); err == nil {
		if err = json.Unmarshal(bytes, &response); err == nil {
			if response.Error == nil {
				return response, nil
			}

			err = response.Error.err()
		}
	} else {
		var res CommonResponse
		if e := json.Unmarshal(bytes, &res); e == nil {
			err = fmt.Errorf("%s: %s", err, res.Error.err())
		}
	}

	return ChatCompletion{}, err
}
