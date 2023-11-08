package openai

// https://platform.openai.com/docs/api-reference/embeddings

import (
	"encoding/json"
	"fmt"
)

// Embeddings struct for response
type Embeddings struct {
	CommonResponse

	Data  []Embedding `json:"data"`
	Model string      `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Embedding struct for Embeddings struct
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingEncodingFormat type for constants
type EmbeddingEncodingFormat string

const (
	EmbeddingEncodingFormatFloat  EmbeddingEncodingFormat = "float"
	EmbeddingEncodingFormatBase64 EmbeddingEncodingFormat = "base64"
)

// EmbeddingOptions for creating embedding
type EmbeddingOptions map[string]any

// SetEncodingFormat sets the `encoding_format` parameter of create embedding request.
//
// https://platform.openai.com/docs/api-reference/embeddings/create#embeddings-create-encoding_format
func (o EmbeddingOptions) SetEncodingFormat(format EmbeddingEncodingFormat) EmbeddingOptions {
	o["encoding_format"] = format
	return o
}

// SetUser sets the `user` parameter of create embedding request.
//
// https://platform.openai.com/docs/api-reference/embeddings/create#embeddings/create-user
func (o EmbeddingOptions) SetUser(user string) EmbeddingOptions {
	o["user"] = user
	return o
}

// CreateEmbedding creates an embedding with given input.
//
// https://platform.openai.com/docs/api-reference/embeddings/create
func (c *Client) CreateEmbedding(model string, input any, options EmbeddingOptions) (response Embeddings, err error) {
	if options == nil {
		options = EmbeddingOptions{}
	}
	options["model"] = model
	options["input"] = input

	var bytes []byte
	if bytes, err = c.post("v1/embeddings", options); err == nil {
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

	return Embeddings{}, err
}
