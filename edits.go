package openai

// https://platform.openai.com/docs/api-reference/edits

import (
	"encoding/json"
	"fmt"
)

// Edit struct for edit response
//
// https://platform.openai.com/docs/api-reference/edits/create
type Edit struct {
	CommonResponse

	Created int64        `json:"created"`
	Choices []EditChoice `json:"choices"`
	Usage   Usage        `json:"usage"`
}

// EditChoice struct for Edit struct
type EditChoice struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

// EditOptions for creating edit
type EditOptions map[string]any

// SetInput sets the `input` parameter of edit request.
//
// https://platform.openai.com/docs/api-reference/edits/create#edits/create-input
func (o EditOptions) SetInput(input string) EditOptions {
	o["input"] = input
	return o
}

// SetN sets the `n` parameter of edit request.
//
// https://platform.openai.com/docs/api-reference/edits/create#edits/create-n
func (o EditOptions) SetN(n int) EditOptions {
	o["n"] = n
	return o
}

// SetTemperature sets the `temperature` parameter of edit request.
//
// https://platform.openai.com/docs/api-reference/edits/create#edits/create-temperature
func (o EditOptions) SetTemperature(temperature float64) EditOptions {
	o["temperature"] = temperature
	return o
}

// SetTopP sets the `top_p` parameter of edit request.
//
// https://platform.openai.com/docs/api-reference/edits/create#edits/create-top_p
func (o EditOptions) SetTopP(topP float64) EditOptions {
	o["top_p"] = topP
	return o
}

// CreateEdit creates an edit for given things.
//
// (DEPRECATED)
//
// https://platform.openai.com/docs/api-reference/edits/create
func (c *Client) CreateEdit(model, instruction string, options EditOptions) (response Edit, err error) {
	if options == nil {
		options = EditOptions{}
	}
	options["model"] = model
	options["instruction"] = instruction

	var bytes []byte
	if bytes, err = c.post("v1/edits", options); err == nil {
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

	return Edit{}, err
}
