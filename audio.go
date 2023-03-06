package openai

import (
	"encoding/json"
	"fmt"
)

// https://platform.openai.com/docs/api-reference/audio

// Transcription struct for response
type Transcription struct {
	CommonResponse

	JSON        *string `json:"json,omitempty"`
	Text        *string `json:"text,omitempty"`
	SRT         *string `json:"srt,omitempty"`
	VerboseJSON *string `json:"verbose_json,omitempty"`
	VTT         *string `json:"vtt,omitempty"`
}

// TranscriptionResponseFormat type for constants
type TranscriptionResponseFormat string

const (
	TranscriptionResponseFormatJSON        TranscriptionResponseFormat = "json"
	TranscriptionResponseFormatText        TranscriptionResponseFormat = "text"
	TranscriptionResponseFormatSRT         TranscriptionResponseFormat = "srt"
	TranscriptionResponseFormatVerboseJSON TranscriptionResponseFormat = "verbose_json"
	TranscriptionResponseFormatVTT         TranscriptionResponseFormat = "vtt"
)

// TranscriptionOptions for creating transcription
type TranscriptionOptions map[string]any

// SetPrompt sets the `prompt` parameter of transcription request.
//
// https://platform.openai.com/docs/api-reference/audio/create#audio/create-prompt
func (o TranscriptionOptions) SetPrompt(prompt string) TranscriptionOptions {
	o["prompt"] = prompt
	return o
}

// SetResponseFormat sets the `response_format` parameter of transcription request.
//
// https://platform.openai.com/docs/api-reference/audio/create#audio/create-response_format
func (o TranscriptionOptions) SetResponseFormat(responseFormat TranscriptionResponseFormat) TranscriptionOptions {
	o["response_format"] = responseFormat
	return o
}

// SetTemperature sets the `temperature` parameter of transcription request.
//
// https://platform.openai.com/docs/api-reference/audio/create#audio/create-temperature
func (o TranscriptionOptions) SetTemperature(temperature float64) TranscriptionOptions {
	o["temperature"] = temperature
	return o
}

// SetLanguage sets the `language` parameter of transcription request.
//
// https://platform.openai.com/docs/api-reference/audio/create#audio/create-language
func (o TranscriptionOptions) SetLanguage(language string) TranscriptionOptions {
	o["language"] = language
	return o
}

// CreateTranscription transcribes given audio file into the input language.
//
// https://platform.openai.com/docs/api-reference/audio/create
func (c *Client) CreateTranscription(file FileParam, model string, options TranscriptionOptions) (response Transcription, err error) {
	if options == nil {
		options = TranscriptionOptions{}
	}
	options["file"] = file
	options["model"] = model

	var bytes []byte
	if bytes, err = c.post("v1/audio/transcriptions", options); err == nil {
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

	return Transcription{}, err
}

// Transcription struct for response
type Translation Transcription

// TransclationResponseFormat type for constants
type TranslationResponseFormat TranscriptionResponseFormat

// TranslationOptions for creating transcription
type TranslationOptions map[string]any

// SetPrompt sets the `prompt` parameter of translation request.
//
// https://platform.openai.com/docs/api-reference/audio/create#audio/create-prompt
func (o TranslationOptions) SetPrompt(prompt string) TranslationOptions {
	o["prompt"] = prompt
	return o
}

// SetResponseFormat sets the `response_format` parameter of translation request.
//
// https://platform.openai.com/docs/api-reference/audio/create#audio/create-response_format
func (o TranslationOptions) SetResponseFormat(responseFormat TranslationResponseFormat) TranslationOptions {
	o["response_format"] = responseFormat
	return o
}

// SetTemperature sets the `temperature` parameter of translation request.
//
// https://platform.openai.com/docs/api-reference/audio/create#audio/create-temperature
func (o TranslationOptions) SetTemperature(temperature float64) TranslationOptions {
	o["temperature"] = temperature
	return o
}

// CreateTranslation translates given audio file into English.
//
// https://platform.openai.com/docs/api-reference/audio/create
func (c *Client) CreateTranslation(file FileParam, model string, options TranslationOptions) (response Translation, err error) {
	if options == nil {
		options = TranslationOptions{}
	}
	options["file"] = file
	options["model"] = model

	var bytes []byte
	if bytes, err = c.post("v1/audio/translations", options); err == nil {
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

	return Translation{}, err
}
