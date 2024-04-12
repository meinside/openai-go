package openai

// https://platform.openai.com/docs/api-reference/chat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// ChatMessageRole type for constants
type ChatMessageRole string

const (
	ChatMessageRoleSystem    ChatMessageRole = "system"
	ChatMessageRoleUser      ChatMessageRole = "user"
	ChatMessageRoleAssistant ChatMessageRole = "assistant"
	ChatMessageRoleTool      ChatMessageRole = "tool"
)

// ToolCall struct
type ToolCall struct {
	Index	 int              `json:"index"`
	ID       string           `json:"id"`
	Type     string           `json:"type"` // == 'function'
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction struct for ToolCall
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ArgumentsParsed returns the parsed map from ToolCall.
func (c ToolCall) ArgumentsParsed() (result map[string]any, err error) {
	if c.Function.Arguments != "" {
		err = json.Unmarshal([]byte(c.Function.Arguments), &result)
	}

	return result, err
}

// ArgumentsInto parses the generated arguments into a given interface
func (c ToolCall) ArgumentsInto(out any) (err error) {
	if c.Function.Arguments == "" {
		err = fmt.Errorf("parse failed: `arguments` is empty")
	} else {
		err = json.Unmarshal([]byte(c.Function.Arguments), &out)
	}

	return err
}

// ChatMessageContent struct
type ChatMessageContent struct {
	Type string `json:"type"`

	Text     *string `json:"text,omitempty"`
	ImageURL any     `json:"image_url,omitempty"`
}

// NewChatMessageContentWithText returns a ChatMessageContent struct with given `text`.
func NewChatMessageContentWithText(text string) ChatMessageContent {
	return ChatMessageContent{
		Type: "text",
		Text: &text,
	}
}

// NewChatMessageContentWithImageURL returns a ChatMessageContent struct with given `url`.
func NewChatMessageContentWithImageURL(url string) ChatMessageContent {
	return ChatMessageContent{
		Type:     "image_url",
		ImageURL: &url,
	}
}

// converts given bytes array to base64-encoded data URL
func bytesToDataURL(bytes []byte) string {
	return fmt.Sprintf("data:%s;base64,%s", http.DetectContentType(bytes), base64.StdEncoding.EncodeToString(bytes))
}

// NewChatMessageContentWithBytes returns a ChatMessageContent struct with given `bytes`.
func NewChatMessageContentWithBytes(bytes []byte) ChatMessageContent {
	return ChatMessageContent{
		Type: "image_url",
		ImageURL: map[string]string{
			"url": bytesToDataURL(bytes),
		},
	}
}

// NewChatMessageContentWithFileParam returns a ChatMessageContent struct with given `file`.
func NewChatMessageContentWithFileParam(file FileParam) ChatMessageContent {
	return NewChatMessageContentWithBytes(file.bs)
}

// ChatMessage struct for chat completion
//
// https://platform.openai.com/docs/guides/chat/introduction
type ChatMessage struct {
	Role    ChatMessageRole `json:"role"`
	Content any             `json:"content,omitempty"` // NOTE: string | []ChatMessageContent

	// for function call
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // when role == 'assistant'
	ToolCallID *string    `json:"tool_call_id,omitempty"` // when role == 'tool'
}

// ContentString tries to return the `content` value as a string.
func (m ChatMessage) ContentString() (string, error) {
	if m.Content != nil {
		if str, ok := m.Content.(string); ok {
			return str, nil
		} else if str, ok := m.Content.(*string); ok { // FIXME: really needed?
			return *str, nil
		}

		return "", fmt.Errorf("returned `content` is not a string")
	}

	return "", fmt.Errorf("returned `content` is nil, cannot return as a string")
}

// ContentArray tries to return the `content` value as a content array.
func (m ChatMessage) ContentArray() ([]ChatMessageContent, error) {
	if m.Content != nil {
		if arr, ok := m.Content.([]ChatMessageContent); ok {
			return arr, nil
		}

		return nil, fmt.Errorf("returned `content` is not a content array")
	}

	return nil, fmt.Errorf("returned `content` is nil, cannot return as a content array")
}

// ToolFunctionParameters type
type ToolFunctionParameters map[string]any

// ChatCompletionTool struct for chat completion function
//
// https://platform.openai.com/docs/api-reference/chat/create#chat-create-tools
type ChatCompletionTool struct {
	Type     string                     `json:"type"` // == 'function'
	Function ChatCompletionToolFunction `json:"function"`
}

// ChatCompletionToolFunction struct
type ChatCompletionToolFunction struct {
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Parameters  ToolFunctionParameters `json:"parameters"`
}

// NewChatCompletionTool returns a ChatCompletionTool.
func NewChatCompletionTool(name, description string, parameters ToolFunctionParameters) ChatCompletionTool {
	tool := ChatCompletionTool{
		Type: "function",
		Function: ChatCompletionToolFunction{
			Name:        name,
			Description: &description,
			Parameters:  parameters,
		},
	}

	return tool
}

// NewToolFunctionParameters returns an empty ToolFunctionParameters.
func NewToolFunctionParameters() ToolFunctionParameters {
	return ToolFunctionParameters{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
}

// AddPropertyWithDescription adds/overwrites a property in chat completion function parameters with a description.
func (p ToolFunctionParameters) AddPropertyWithDescription(name, typ3, description string) ToolFunctionParameters {
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

// AddArrayPropertyWithDescription adds/overwrites an array property in chat completion function parameters with a description.
func (p ToolFunctionParameters) AddArrayPropertyWithDescription(name, elemType, description string) ToolFunctionParameters {
	if properties, exists := p["properties"]; exists {
		ps := properties.(map[string]any)
		ps[name] = map[string]any{
			"type":        "array",
			"description": description,
			"items": map[string]any{
				"type": elemType,
			},
		}
		p["properties"] = ps
	}
	return p
}

// AddPropertyWithEnums adds/overwrites a property in chat completion function parameters with enums.
func (p ToolFunctionParameters) AddPropertyWithEnums(name, typ3, description string, enums []string) ToolFunctionParameters {
	if properties, exists := p["properties"]; exists {
		ps := properties.(map[string]any)
		ps[name] = map[string]any{
			"type":        typ3,
			"description": description,
			"enum":        enums,
		}
		p["properties"] = ps
	}
	return p
}

// SetRequiredParameters sets/overwrites required parameter names for chat completion function parameters.
func (p ToolFunctionParameters) SetRequiredParameters(names []string) ToolFunctionParameters {
	p["required"] = names
	return p
}

// ChatCompletionToolChoiceMode type
//
// https://platform.openai.com/docs/api-reference/chat/create#chat-create-tool_choice
type ChatCompletionToolChoiceMode string

// ChatCompletionToolChoiceMode constants
const (
	ChatCompletionToolChoiceNone ChatCompletionToolChoiceMode = "none"
	ChatCompletionToolChoiceAuto ChatCompletionToolChoiceMode = "auto"
)

// NewChatSystemMessage returns a new ChatMessage with system role.
func NewChatSystemMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleSystem,
		Content: &message,
	}
}

// ChatUserMessageContentTypes interface for type constraints in `NewChatUserMessage`
type ChatUserMessageContentTypes interface {
	string | []ChatMessageContent
}

// NewChatUserMessage returns a new ChatMessage with user role.
func NewChatUserMessage[T ChatUserMessageContentTypes](contents T) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleUser,
		Content: contents,
	}
}

// NewChatAssistantMessage returns a new ChatMessage with assistant role.
func NewChatAssistantMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleAssistant,
		Content: &message,
	}
}

// NewChatToolMessage returns a new ChatMesssage with tool role.
func NewChatToolMessage(toolCallID, content string) ChatMessage {
	return ChatMessage{
		Role:       ChatMessageRoleTool,
		Content:    &content,
		ToolCallID: &toolCallID,
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

// ChatCompletionResponseFormat struct for chat completion request
type ChatCompletionResponseFormat struct {
	Type ChatCompletionResponseFormatType `json:"type,omitempty"`
}

// ChatCompletionResponseFormatType type for constants
type ChatCompletionResponseFormatType string

// ChatCompletionResponseFormatType constants
const (
	ChatCompletionResponseFormatTypeText       ChatCompletionResponseFormatType = "text"
	ChatCompletionResponseFormatTypeJSONObject ChatCompletionResponseFormatType = "json_object"
)

// ChatCompletionOptions for creating chat completions
type ChatCompletionOptions map[string]any

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

// SetMaxTokens sets the `max_tokens` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-max_tokens
func (o ChatCompletionOptions) SetMaxTokens(maxTokens int) ChatCompletionOptions {
	o["max_tokens"] = maxTokens
	return o
}

// SetN sets the `n` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-n
func (o ChatCompletionOptions) SetN(n int) ChatCompletionOptions {
	o["n"] = n
	return o
}

// SetPresencePenalty sets the `presence_penalty` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-presence_penalty
func (o ChatCompletionOptions) SetPresencePenalty(presencePenalty float64) ChatCompletionOptions {
	o["presence_penalty"] = presencePenalty
	return o
}

// SetResponseFormat sets the `response_format` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat-create-response_format
func (o ChatCompletionOptions) SetResponseFormat(format ChatCompletionResponseFormat) ChatCompletionOptions {
	o["response_format"] = format
	return o
}

// SetSeed sets the `seed` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat-create-seed
func (o ChatCompletionOptions) SetSeed(seed int64) ChatCompletionOptions {
	o["seed"] = seed
	return o
}

// SetStop sets the `stop` parameter of chat completions.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat/create-stop
func (o ChatCompletionOptions) SetStop(stop any) ChatCompletionOptions {
	o["stop"] = stop
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

// SetTools sets the `tools` parameter of chat completion request.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat-create-tools
func (o ChatCompletionOptions) SetTools(tools []ChatCompletionTool) ChatCompletionOptions {
	o["tools"] = tools
	return o
}

// SetToolChoice sets the `tool_choice` parameter of chat completion request.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat-create-tool_choice
func (o ChatCompletionOptions) SetToolChoice(choice ChatCompletionToolChoiceMode) ChatCompletionOptions {
	o["tool_choice"] = choice
	return o
}

// SetToolChoiceWithName sets the `tool_choice` parameter of chat completion request with given function name.
//
// https://platform.openai.com/docs/api-reference/chat/create#chat-create-tool_choice
func (o ChatCompletionOptions) SetToolChoiceWithName(name string) ChatCompletionOptions {
	o["tool_choice"] = map[string]any{
		"type": "function",
		"function": map[string]any{
			"name": name,
		},
	}
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
