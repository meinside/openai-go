package openai

import (
	"encoding/json"
	"fmt"
)

// https://platform.openai.com/docs/api-reference/moderations

// Moderation struct for response
type Moderation struct {
	CommonResponse

	ID      string           `json:"id"`
	Model   string           `json:"model"`
	Results []Classification `json:"results"`
}

// Classification struct for Moderation strcut
type Classification struct {
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float64 `json:"category_scores"`
	Flagged        bool               `json:"flagged"`
}

// ModerationOptions for creating moderation
type ModerationOptions map[string]any

// SetModel sets the `model` parameter of moderation request.
//
// https://platform.openai.com/docs/api-reference/moderations/create#moderations/create-model
func (o ModerationOptions) SetModel(model string) ModerationOptions {
	o["model"] = model
	return o
}

// CreateModeration classifies given text.
//
// https://platform.openai.com/docs/api-reference/moderations/create
func (c *Client) CreateModeration(input any, options ModerationOptions) (response Moderation, err error) {
	if options == nil {
		options = ModerationOptions{}
	}
	options["input"] = input

	var bytes []byte
	if bytes, err = c.post("v1/moderations", options); err == nil {
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

	return Moderation{}, err
}
