package openai

import (
	"encoding/json"
	"fmt"
)

// https://platform.openai.com/docs/api-reference/images

// GeneratedImages struct for image creation responses
type GeneratedImages struct {
	CommonResponse

	Created int64 `json:"created"`
	Data    []struct {
		URL        *string `json:"url,omitempty"`
		Base64JSON *string `json:"b64_json,omitempty"`
	} `json:"data"`
}

// ImageSize type for constants
type ImageSize string

const (
	ImageSize256x256   ImageSize = "256x256"
	ImageSize512x512   ImageSize = "512x512"
	ImageSize1024x1024 ImageSize = "1024x1024"
)

// ImageResponseFormat type for constants
type ImageResponseFormat string

const (
	IamgeResponseFormatURL        ImageResponseFormat = "url"
	IamgeResponseFormatBase64JSON ImageResponseFormat = "b64_json"
)

// ImageOptions for creating images
type ImageOptions map[string]any

// SetN sets the `n` parameter of image generation request.
//
// https://platform.openai.com/docs/api-reference/images/create#images/create-n
func (o ImageOptions) SetN(n int) ImageOptions {
	o["n"] = n
	return o
}

// SetSize sets the `size` parameter of image generation request.
//
// https://platform.openai.com/docs/api-reference/images/create#images/create-size
func (o ImageOptions) SetSize(size ImageSize) ImageOptions {
	o["size"] = size
	return o
}

// SetResponseFormat sets the `response_format` parameter of image generation request.
//
// https://platform.openai.com/docs/api-reference/images/create#images/create-response_format
func (o ImageOptions) SetResponseFormat(responseFormat ImageResponseFormat) ImageOptions {
	o["response_format"] = responseFormat
	return o
}

// SetUser sets the `user` parameter of image generation request.
//
// https://platform.openai.com/docs/api-reference/images/create#images/create-user
func (o ImageOptions) SetUser(user string) ImageOptions {
	o["user"] = user
	return o
}

// CreateImage creates an image with given prompt.
//
// https://platform.openai.com/docs/api-reference/images/create
func (c *Client) CreateImage(prompt string, options ImageOptions) (response GeneratedImages, err error) {
	if options == nil {
		options = ImageOptions{}
	}
	options["prompt"] = prompt

	var bytes []byte
	if bytes, err = c.post("v1/images/generations", options); err == nil {
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

	return GeneratedImages{}, err
}

// ImageEditOptions for creating image edits
type ImageEditOptions map[string]any

// SetMask sets the `mask` parameter of image edit request.
//
// https://platform.openai.com/docs/api-reference/images/create-edit#images/create-edit-mask
func (o ImageEditOptions) SetMask(mask FileParam) ImageEditOptions {
	o["mask"] = mask
	return o
}

// SetN sets the `n` parameter of image edit request.
//
// https://platform.openai.com/docs/api-reference/images/create-edit#images/create-edit-n
func (o ImageEditOptions) SetN(n int) ImageEditOptions {
	o["n"] = n
	return o
}

// SetSize sets the `size` parameter of image edit request.
//
// https://platform.openai.com/docs/api-reference/images/create-edit#images/create-edit-size
func (o ImageEditOptions) SetSize(size ImageSize) ImageEditOptions {
	o["size"] = size
	return o
}

// SetResponseFormat sets the `response_format` parameter of image edit request.
//
// https://platform.openai.com/docs/api-reference/images/create-edit#images/create-edit-response_format
func (o ImageEditOptions) SetResponseFormat(responseFormat ImageResponseFormat) ImageEditOptions {
	o["response_format"] = responseFormat
	return o
}

// SetUser sets the `user` parameter of image edit request.
//
// https://platform.openai.com/docs/api-reference/images/create-edit#images/create-edit-user
func (o ImageEditOptions) SetUser(user string) ImageEditOptions {
	o["user"] = user
	return o
}

// CreateImageEdit creates an edited or extended image with given file and prompt.
//
// https://platform.openai.com/docs/api-reference/images/create-edit
func (c *Client) CreateImageEdit(image FileParam, prompt string, options ImageEditOptions) (response GeneratedImages, err error) {
	if options == nil {
		options = ImageEditOptions{}
	}
	options["image"] = image
	options["prompt"] = prompt

	var bytes []byte
	if bytes, err = c.post("v1/images/edits", options); err == nil {
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

	return GeneratedImages{}, err
}

// ImageVariationOptions for creating image variations
type ImageVariationOptions map[string]any

// SetN sets the `n` parameter of image variation request.
//
// https://platform.openai.com/docs/api-reference/images/create-variation#images/create-variation-n
func (o ImageVariationOptions) SetN(n int) ImageVariationOptions {
	o["n"] = n
	return o
}

// SetSize sets the `size` parameter of image variation request.
//
// https://platform.openai.com/docs/api-reference/images/create-variation#images/create-variation-size
func (o ImageVariationOptions) SetSize(size ImageSize) ImageVariationOptions {
	o["size"] = size
	return o
}

// SetResponseFormat sets the `response_format` parameter of image variation request.
//
// https://platform.openai.com/docs/api-reference/images/create-variation#images/create-variation-response_format
func (o ImageVariationOptions) SetResponseFormat(responseFormat ImageResponseFormat) ImageVariationOptions {
	o["response_format"] = responseFormat
	return o
}

// SetUser sets the `user` parameter of image variation request.
//
// https://platform.openai.com/docs/api-reference/images/create-variation#images/create-variation-user
func (o ImageVariationOptions) SetUser(user string) ImageVariationOptions {
	o["user"] = user
	return o
}

// CreateImageVariation creates a variation of a given image.
//
// https://platform.openai.com/docs/api-reference/images/create-variation
func (c *Client) CreateImageVariation(image FileParam, options ImageVariationOptions) (response GeneratedImages, err error) {
	if options == nil {
		options = ImageVariationOptions{}
	}
	options["image"] = image

	var bytes []byte
	if bytes, err = c.post("v1/images/variations", options); err == nil {
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

	return GeneratedImages{}, err
}
