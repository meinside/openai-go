package openai

import (
	"encoding/json"
	"fmt"
)

// https://platform.openai.com/docs/api-reference/fine-tunes

// FineTune struct for response
type FineTune struct {
	CommonResponse

	ID             string          `json:"id"`
	Model          string          `json:"model"`
	CreatedAt      int64           `json:"created_at"`
	Events         []FineTuneEvent `json:"events,omitempty"`
	FineTunedModel *string         `json:"fine_tuned_model,omitempty"`
	HyperParams    struct {
		BatchSize              int     `json:"batch_size"`
		LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
		NEpochs                int     `json:"n_epochs"`
		PromptLossWeight       float64 `json:"prompt_loss_weight"`
	} `json:"hyperparams"`
	OrganizationID  string `json:"organization_id"`
	ResultFiles     []File `json:"result_files"`
	Status          string `json:"status"`
	ValidationFiles []File `json:"validation_files"`
	TrainingFiles   []File `json:"training_files"`
	UpdatedAt       int64  `json:"updated_at"`
}

// FileTuneEvent struct for response
type FineTuneEvent struct {
	CommonResponse

	CreatedAt int64  `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// FineTuneOptions for creating moderation
type FineTuneOptions map[string]any

// SetValidationFile sets the `validation_file` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-validation_file
func (o FineTuneOptions) SetValidationFile(validationFileID string) FineTuneOptions {
	o["validation_file"] = validationFileID
	return o
}

// SetModel sets the `model` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-model
func (o FineTuneOptions) SetModel(model string) FineTuneOptions {
	o["model"] = model
	return o
}

// SetNEpochs sets the `n_epochs` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-n_epochs
func (o FineTuneOptions) SetNEpochs(nEpochs int) FineTuneOptions {
	o["n_epochs"] = nEpochs
	return o
}

// SetBatchSize sets the `batch_size` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-batch_size
func (o FineTuneOptions) SetBatchSize(batchSize int) FineTuneOptions {
	o["batch_size"] = batchSize
	return o
}

// SetLearningRateMultiplier sets the `learning_rate_multiplier` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-learning_rate_multiplier
func (o FineTuneOptions) SetLearningRateMultiplier(learningRateMultiplier float64) FineTuneOptions {
	o["learning_rate_multiplier"] = learningRateMultiplier
	return o
}

// SetPromptLossWeight sets the `prompt_loss_weight` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-prompt_loss_weight
func (o FineTuneOptions) SetPromptLossWeight(promptLossWeight float64) FineTuneOptions {
	o["prompt_loss_weight"] = promptLossWeight
	return o
}

// SetComputeClassificationMetrics sets the `compute_classification_metrics` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-compute_classification_metrics
func (o FineTuneOptions) SetComputeClassificationMetrics(computeClassificationMetrics bool) FineTuneOptions {
	o["compute_classification_metrics"] = computeClassificationMetrics
	return o
}

// SetClassificaitonNClasses sets the `classification_n_classes` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-classification_n_classes
func (o FineTuneOptions) SetClassificationNClasses(classificationNClasses int) FineTuneOptions {
	o["classification_n_classes"] = classificationNClasses
	return o
}

// SetClassificationPositiveClass sets the `classification_positive_class` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-classification_positive_class
func (o FineTuneOptions) SetClassificationPositiveClass(classificationPositiveClass string) FineTuneOptions {
	o["classification_positive_class"] = classificationPositiveClass
	return o
}

// SetClassificationBetas sets the `classification_betas` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-classification_betas
func (o FineTuneOptions) SetClassificationBetas(classificationBetas []float64) FineTuneOptions {
	o["classification_betas"] = classificationBetas
	return o
}

// SetSuffix sets the `suffix` parameter of fine-tune request.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-suffix
func (o FineTuneOptions) SetSuffix(suffix string) FineTuneOptions {
	o["suffix"] = suffix
	return o
}

// CreateFineTune creates a job that fine-tunes a specified model from a given dataset.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/create
func (c *Client) CreateFineTune(trainingFileID string, options FineTuneOptions) (response FineTune, err error) {
	if options == nil {
		options = FineTuneOptions{}
	}
	options["training_file"] = trainingFileID

	var bytes []byte
	if bytes, err = c.post("v1/fine-tunes", options); err == nil {
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

	return FineTune{}, err
}

// FineTunes struct for response
type FineTunes struct {
	CommonResponse

	Data []FineTune `json:"data"`
}

// ListFineTunes lists the organization's fine-tuning jobs.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/list
func (c *Client) ListFineTunes() (response FineTunes, err error) {
	var bytes []byte
	if bytes, err = c.get("v1/fine-tunes", nil); err == nil {
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

	return FineTunes{}, err
}

// RetrieveFineTune gets info about specified fine-tune job.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/retrieve
func (c *Client) RetrieveFineTune(fineTuneID string) (response FineTune, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/fine-tunes/%s", fineTuneID), nil); err == nil {
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

	return FineTune{}, err
}

// CancelFineTune immediately cancels a fine-tune job.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/cancel
func (c *Client) CancelFineTune(fineTuneID string) (response FineTune, err error) {
	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/fine-tunes/%s/cancel", fineTuneID), nil); err == nil {
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

	return FineTune{}, err
}

// FineTuneEvents struct for response
type FineTuneEvents struct {
	CommonResponse

	Data []FineTuneEvent `json:"data"`
}

// FineTuneEventsOptions for listing fine-tune events
type FineTuneEventsOptions map[string]any

// SetStream sets the `stream` parameter of fine-tune request.
//
// NOTE: (not implemented) https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format
//
// https://platform.openai.com/docs/api-reference/fine-tunes/events#fine-tunes/events-stream
func (o FineTuneEventsOptions) SetStream(stream bool) FineTuneEventsOptions {
	o["stream"] = stream
	return o
}

// ListFineTuneEvents gets fine-grained status updates for a fine-tune job.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/events
func (c *Client) ListFineTuneEvents(fineTuneID string, options FineTuneEventsOptions) (response FineTuneEvents, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/fine-tunes/%s/events", fineTuneID), options); err == nil {
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

	return FineTuneEvents{}, err
}

// DeletedFineTuneModel struct for response
type DeletedFineTuneModel struct {
	CommonResponse

	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// DeleteFineTuneModel deletes a fine-tune model.
//
// https://platform.openai.com/docs/api-reference/fine-tunes/delete-model
func (c *Client) DeleteFineTuneModel(model string) (response DeletedFineTuneModel, err error) {
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

	return DeletedFineTuneModel{}, err
}
