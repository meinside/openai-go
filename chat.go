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
)

// ChatMessage struct for chat completion
//
// https://platform.openai.com/docs/guides/chat/introduction
type ChatMessage struct {
	Role    ChatMessageRole `json:"role"`
	Content string          `json:"content"`
}

// NewChatSystemMessage returns a new ChatMessage with system role.
func NewChatSystemMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleSystem,
		Content: message,
	}
}

// NewChatUserMessage returns a new ChatMessage with user role.
func NewChatUserMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleUser,
		Content: message,
	}
}

// NewChatAssistantMessage returns a new ChatMessage with assistant role.
func NewChatAssistantMessage(message string) ChatMessage {
	return ChatMessage{
		Role:    ChatMessageRoleAssistant,
		Content: message,
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
