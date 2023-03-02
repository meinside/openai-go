package openai

// https://platform.openai.com/docs/api-reference/completions

import (
	"encoding/json"
	"fmt"
)

// CompletionOptions for creating completions
type CompletionOptions map[string]any

// SetPrompt sets the `prompt` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-prompt
func (o CompletionOptions) SetPrompt(prompt any) CompletionOptions {
	o["prompt"] = prompt
	return o
}

// SetSuffix sets the `suffix` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-suffix
func (o CompletionOptions) SetSuffix(suffix string) CompletionOptions {
	o["suffix"] = suffix
	return o
}

// SetMaxTokens sets the `max_tokens` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-max_tokens
func (o CompletionOptions) SetMaxTokens(maxTokens int) CompletionOptions {
	o["max_tokens"] = maxTokens
	return o
}

// SetTemperature sets the `temperature` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-temperature
func (o CompletionOptions) SetTemperature(temperature float64) CompletionOptions {
	o["temperature"] = temperature
	return o
}

// SetTopP sets the `top_p` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-top_p
func (o CompletionOptions) SetTopP(topP float64) CompletionOptions {
	o["top_p"] = topP
	return o
}

// SetN sets the `n` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-n
func (o CompletionOptions) SetN(n int) CompletionOptions {
	o["n"] = n
	return o
}

// SetStream sets the `stream` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-stream
func (o CompletionOptions) SetStream(stream bool) CompletionOptions {
	o["stream"] = stream
	return o
}

// SetLogProbabilities sets the `logprobs` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-logprobs
func (o CompletionOptions) SetLogProbabilities(logprobs int) CompletionOptions {
	o["logprobs"] = logprobs
	return o
}

// SetEcho sets the `echo` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-echo
func (o CompletionOptions) SetEcho(echo bool) CompletionOptions {
	o["echo"] = echo
	return o
}

// SetStop sets the `stop` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-stop
func (o CompletionOptions) SetStop(stop any) CompletionOptions {
	o["stop"] = stop
	return o
}

// SetPresencePenalty sets the `presence_penalty` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-presence_penalty
func (o CompletionOptions) SetPresencePenalty(presencePenalty float64) CompletionOptions {
	o["presence_penalty"] = presencePenalty
	return o
}

// SetFrequencyPenalty sets the `frequency_penalty` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-frequency_penalty
func (o CompletionOptions) SetFrequencyPenalty(frequencyPenalty float64) CompletionOptions {
	o["frequency_penalty"] = frequencyPenalty
	return o
}

// SetBestOf sets the `best_of` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-best_of
func (o CompletionOptions) SetBestOf(bestOf int) CompletionOptions {
	o["best_of"] = bestOf
	return o
}

// SetLogitBias sets the `logit_bias` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-logit_bias
func (o CompletionOptions) SetLogitBias(logitBias map[string]any) CompletionOptions {
	o["logit_bias"] = logitBias
	return o
}

// SetUser sets the `user` parameter of completion request.
//
// https://platform.openai.com/docs/api-reference/completions/create#completions/create-user
func (o CompletionOptions) SetUser(user string) CompletionOptions {
	o["user"] = user
	return o
}

// CompletionChoice struct for completion response
type CompletionChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     *int   `json:"logprobs,omitempty"`
	FinishReason string `json:"finish_reason"`
}

// Completion struct for response
type Completion struct {
	CommonResponse

	ID      string             `json:"id"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// CreateCompletion creates a completion.
//
// https://platform.openai.com/docs/api-reference/completions/create
func (c *Client) CreateCompletion(model string, options CompletionOptions) (response Completion, err error) {
	if options == nil {
		options = CompletionOptions{}
	}
	options["model"] = model

	var bytes []byte
	if bytes, err = c.post("v1/completions", options); err == nil {
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

	return Completion{}, err
}
