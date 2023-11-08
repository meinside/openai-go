package openai

import (
	"encoding/json"
	"fmt"
)

// https://platform.openai.com/docs/api-reference/fine-tuning

// FineTuningJobs struct
type FineTuningJobs struct {
	CommonResponse

	Data    []FineTuningJob `json:"data"`
	HasMore bool            `json:"has_more"`
}

// FineTuningJob struct
type FineTuningJob struct {
	CommonResponse

	ID              string                    `json:"id"`
	CreatedAt       int64                     `json:"created_at"`
	FinishedAt      int64                     `json:"finished_at"`
	Model           string                    `json:"model"`
	FineTunedModel  *string                   `json:"fine_tuned_model,omitempty"`
	OrganizationID  string                    `json:"organization_id"`
	Status          FineTuningJobStatus       `json:"status"`
	Hyperparameters FineTuningHyperparameters `json:"hyperparameters"`
	TrainingFile    string                    `json:"training_file"`
	ValidationFile  *string                   `json:"validation_file,omitempty"`
	ResultFiles     []string                  `json:"result_files"`
	TrainedTokens   int                       `json:"trained_tokens"`
}

// FineTuningJobStatus type
type FineTuningJobStatus string

// FineTuningJobStatus constants
const (
	FineTuningJobStatusCreated   FineTuningJobStatus = "created"
	FineTuningJobStatusPending   FineTuningJobStatus = "pending"
	FineTuningJobStatusRunning   FineTuningJobStatus = "running"
	FineTuningJobStatusSucceeded FineTuningJobStatus = "succeeded"
	FineTuningJobStatusFailed    FineTuningJobStatus = "failed"
	FineTuningJobStatusCancelled FineTuningJobStatus = "cancelled"
)

// FineTuningHyperparameters struct
type FineTuningHyperparameters struct {
	NEpochs any `json:"n_epochs"` // string("Auto") or int
}

// FineTuningJobsOptions for listing fine-tuning jobs
type FineTuningJobsOptions map[string]any

// SetAfter sets the `after` parameter of fine-tuning jobs listing request.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/list#fine-tuning-list-after
func (o FineTuningJobsOptions) SetAfter(after string) FineTuningJobsOptions {
	o["after"] = after
	return o
}

// SetLimit sets the `limit` parameter of fine-tuning jobs listing request.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/list#fine-tuning-list-limit
func (o FineTuningJobsOptions) SetLimit(limit int) FineTuningJobsOptions {
	o["limit"] = limit
	return o
}

// FineTuningJobOptions for retrieving fine-tuning jobs
type FineTuningJobOptions map[string]any

// SetValidationFile sets the `validation_file` parameter of fine-tuning job request.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/create#validation_file
func (o FineTuningJobOptions) SetValidationFile(validationFileID string) FineTuningJobOptions {
	o["validation_file"] = validationFileID
	return o
}

// SetHyperparameters sets the `hyperparameters` parameter of fine-tuning job request.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/create#model
func (o FineTuningJobOptions) SetHyperparameters(hyperparameters FineTuningHyperparameters) FineTuningJobOptions {
	o["hyperparameters"] = hyperparameters
	return o
}

// SetSuffix sets the `suffix` parameter of fine-tuning job request.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/create#suffix
func (o FineTuningJobOptions) SetSuffix(suffix string) FineTuningJobOptions {
	o["suffix"] = suffix
	return o
}

// CreateFineTuningJob creates a job that fine-tunes a specified model from given data
//
// https://platform.openai.com/docs/api-reference/fine-tuning/create
func (c *Client) CreateFineTuningJob(trainingFileID, model string, options FineTuningJobOptions) (response FineTuningJob, err error) {
	if options == nil {
		options = FineTuningJobOptions{}
	}
	options["training_file"] = trainingFileID
	options["model"] = model

	var bytes []byte
	if bytes, err = c.post("v1/fine_tuning/jobs", options); err == nil {
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

	return FineTuningJob{}, err
}

// ListFineTuningJobs lists your organization's fine-tuning jobs.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/list
func (c *Client) ListFineTuningJobs(options FineTuningJobsOptions) (response FineTuningJobs, err error) {
	if options == nil {
		options = FineTuningJobsOptions{}
	}

	var bytes []byte
	if bytes, err = c.get("v1/fine_tuning/jobs", options); err == nil {
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

	return FineTuningJobs{}, err
}

// RetrieveFineTuningJob retrieves a fine-tuning job.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/retrieve
func (c *Client) RetrieveFineTuningJob(fineTuningJobID string) (response FineTuningJob, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/fine_tuning/jobs/%s", fineTuningJobID), nil); err == nil {
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

	return FineTuningJob{}, err
}

// CancelFineTuningJob cancels a fine-tuning job.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/cancel
func (c *Client) CancelFineTuningJob(fineTuningJobID string) (response FineTuningJob, err error) {
	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/fine_tuning/jobs/%s/cancel", fineTuningJobID), nil); err == nil {
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

	return FineTuningJob{}, err
}

// FineTuningJobEvents struct
type FineTuningJobEvents struct {
	CommonResponse

	Data []FineTuningJobEvent `json:"data"`

	HasMore bool `json:"has_more"`
}

// FineTuningJobEvent struct
type FineTuningJobEvent struct {
	CommonResponse

	ID        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Type      string `json:"type"`
}

// FineTuningJobEventsOptions for listing fine-tuning job events
type FineTuningJobEventsOptions map[string]any

// SetAfter sets the `after` parameter of fine-tuning events request.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/list-events#after
func (o FineTuningJobEventsOptions) SetAfter(fineTuningJobEventID string) FineTuningJobEventsOptions {
	o["after"] = fineTuningJobEventID
	return o
}

// SetLimit sets the `limit` parameter of fine-tuning events request.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/list-events#limit
func (o FineTuningJobEventsOptions) SetLimit(limit int) FineTuningJobEventsOptions {
	o["limit"] = limit
	return o
}

// ListFineTuningJobEvents lists status updates for a given fine-tuning job.
//
// https://platform.openai.com/docs/api-reference/fine-tuning/list-events
func (c *Client) ListFineTuningJobEvents(fineTuningJobID string, options FineTuningJobEventsOptions) (response FineTuningJobEvents, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/fine_tuning/jobs/%s/events", fineTuningJobID), options); err == nil {
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

	return FineTuningJobEvents{}, err
}
