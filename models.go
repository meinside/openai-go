package openai

// https://platform.openai.com/docs/api-reference/models

import (
	"encoding/json"
	"fmt"
)

// Permission struct
type Permission struct {
	ID                 string  `json:"id"`
	Object             string  `json:"object"`
	Created            int64   `json:"created"`
	AllowCreateEngine  bool    `json:"allow_create_engine"`
	AllowSampling      bool    `json:"allow_sampling"`
	AllowLogProbs      bool    `json:"allow_logprobs"`
	AllowSearchIndices bool    `json:"allow_search_indices"`
	AllowView          bool    `json:"allow_view"`
	AllowFineTuning    bool    `json:"allow_fine_tuning"`
	Organization       string  `json:"organization"`
	Group              *string `json:"group,omitempty"`
	IsBlocking         bool    `json:"is_blocking"`
}

// Model struct
type Model struct {
	CommonResponse

	ID         string       `json:"id"`
	Created    int64        `json:"created"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
	Root       string       `json:"root"`
	Parent     *string      `json:"parent,omitempty"`
}

// ModelList struct for API response
type ModelsList struct {
	CommonResponse

	Data []Model `json:"data"`
}

// ModelDeletionStatus struct for API response
type ModelDeletionStatus struct {
	CommonResponse

	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// ListModels lists currently available models.
//
// https://platform.openai.com/docs/api-reference/models/list
func (c *Client) LitModels() (response ModelsList, err error) {
	var bytes []byte
	if bytes, err = c.get("v1/models", nil); err == nil {
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

	return ModelsList{}, err
}

// RetrieveModel retrieves a model instance.
//
// https://platform.openai.com/docs/api-reference/models/retrieve
func (c *Client) RetrieveModel(id string) (response Model, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/models/%s", id), nil); err == nil {
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

	return Model{}, err
}

// DeleteFineTuneModel deletes a fine-tuned model.
//
// https://platform.openai.com/docs/api-reference/models/delete
func (c *Client) DeleteFineTuneModel(model string) (response ModelDeletionStatus, err error) {
	var bytes []byte
	if bytes, err = c.delete(fmt.Sprintf("v1/models/%s", model), nil); err == nil {
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

	return ModelDeletionStatus{}, err
}
