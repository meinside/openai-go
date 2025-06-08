package openai

import (
	"encoding/json"
	"fmt"
)

// NOTE: In Beta

// https://platform.openai.com/docs/api-reference/runs

// RunStatus type for constants
type RunStatus string

// RunStatus constants
const (
	RunStatusQueued         RunStatus = "queued"
	RunStatusInProgress     RunStatus = "in_progress"
	RunStatusRequiresAction RunStatus = "requires_action"
	RunStatusCanceling      RunStatus = "cancelling"
	RunStatusCanceled       RunStatus = "cancelled"
	RunStatusFailed         RunStatus = "failed"
	RunStatusCompleted      RunStatus = "completed"
	RunStatusExpired        RunStatus = "expired"
)

// https://platform.openai.com/docs/api-reference/runs/object
type Run struct {
	CommonResponse

	ID             string            `json:"id"`
	CreatedAt      int               `json:"created_at"`
	ThreadID       string            `json:"thread_id"`
	AssistantID    string            `json:"assistant_id"`
	Status         RunStatus         `json:"status"`
	RequiredAction *RunAction        `json:"required_action,omitempty"`
	LastError      *RunError         `json:"last_error,omitempty"`
	ExpiresAt      int               `json:"expires_at"`
	StartedAt      *int              `json:"started_at,omitempty"`
	CancelledAt    *int              `json:"cancelled_at,omitempty"`
	FailedAt       *int              `json:"failed_at,omitempty"`
	CompletedAt    *int              `json:"completed_at,omitempty"`
	Model          string            `json:"model"`
	Instructions   string            `json:"instructions"`
	Tools          []Tool            `json:"tools"`
	FileIDs        []string          `json:"file_ids"`
	Metadata       map[string]string `json:"metadata"`
}

// RunAction struct for Run struct
type RunAction struct {
	Type              string `json:"type"` // == 'submit_tool_outputs'
	SubmitToolOutputs struct {
		ToolCalls []ToolCall `json:"tool_calls"`
	} `json:"submit_tool_outputs"`
}

// RunErrorCode type for constants
type RunErrorCode string

// RunErrorCode constants
const (
	RunErrorCodeServerError   RunErrorCode = "server_error"
	RunErrorRateLimitExceeded RunErrorCode = "rate_limit_exceeded"
)

// RunError struct for Run struct
type RunError struct {
	Code    RunErrorCode `json:"code"`
	Message string       `json:"message"`
}

// CreateRunOptions for creating run
type CreateRunOptions map[string]any

// SetModel sets the `model` parameter of CreateRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createRun#runs-createrun-model
func (o CreateRunOptions) SetModel(model string) CreateRunOptions {
	o["model"] = model
	return o
}

// SetInstructions sets the `instructions` parameter of CreateRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createRun#runs-createrun-instructions
func (o CreateRunOptions) SetInstructions(instructions string) CreateRunOptions {
	o["instructions"] = instructions
	return o
}

// SetTools sets the `tools` parameter of CreateRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createRun#runs-createrun-tools
func (o CreateRunOptions) SetTools(tools []Tool) CreateRunOptions {
	o["tools"] = tools
	return o
}

// SetMetadata sets the `metadata` parameter of CreateRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createRun#runs-createrun-metadata
func (o CreateRunOptions) SetMetadata(metadata map[string]string) CreateRunOptions {
	o["metadata"] = metadata
	return o
}

// CreateRun creates a run with given `threadID`, `assistantID`, and `options`.
//
// https://platform.openai.com/docs/api-reference/runs/createRun
func (c *Client) CreateRun(threadID, assistantID string, options CreateRunOptions) (response Run, err error) {
	if options == nil {
		options = CreateRunOptions{}
	}
	options["assistant_id"] = assistantID

	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/threads/%s/runs", threadID), options); err == nil {
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

	return Run{}, err
}

// RetrieveRun retrieves a run with given `threadID` and `runID`.
//
// https://platform.openai.com/docs/api-reference/runs/getRun
func (c *Client) RetrieveRun(threadID, runID string) (response Run, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/runs/%s", threadID, runID), nil); err == nil {
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

	return Run{}, err
}

// ModifyRunOptions for modifying run
type ModifyRunOptions map[string]any

// SetMetadata sets the `metadata` parameter of ModifyRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/modifyRun#runs-modifyrun-metadata
func (o ModifyRunOptions) SetMetadata(metadata map[string]string) ModifyRunOptions {
	o["metadata"] = metadata
	return o
}

// ModifyRun modifies a run with given `threadID`, `runID`, and `options`.
//
// https://platform.openai.com/docs/api-reference/runs/modifyRun
func (c *Client) ModifyRun(threadID, runID string, options ModifyRunOptions) (response Run, err error) {
	if options == nil {
		options = ModifyRunOptions{}
	}

	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/threads/%s/runs/%s", threadID, runID), options); err == nil {
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

	return Run{}, err
}

// Runs struct for API response
type Runs struct {
	CommonResponse

	Data    []Run  `json:"data"`
	FirstID string `json:"first_id"`
	LastID  string `json:"last_id"`
	HasMore bool   `json:"has_more"`
}

// ListRunsOptions for listing runs
type ListRunsOptions map[string]any

// SetLimit sets the `limit` parameter of messages' listing request.
//
// https://platform.openai.com/docs/api-reference/runs/listRuns#runs-listruns-limit
func (o ListRunsOptions) SetLimit(limit int) ListRunsOptions {
	o["limit"] = limit
	return o
}

// SetOrder sets the `order` parameter of messages' listing request.
//
// `order` can be one of 'asc' or 'desc'. (default: 'desc')
//
// https://platform.openai.com/docs/api-reference/runs/listRuns#runs-listruns-order
func (o ListRunsOptions) SetOrder(order string) ListRunsOptions {
	o["order"] = order
	return o
}

// SetAfter sets the `after` parameter of messages' listing request.
//
// https://platform.openai.com/docs/api-reference/runs/listRuns#runs-listruns-after
func (o ListRunsOptions) SetAfter(after string) ListRunsOptions {
	o["after"] = after
	return o
}

// SetBefore sets the `before` parameter of messages' listing request.
//
// https://platform.openai.com/docs/api-reference/runs/listRuns#runs-listruns-before
func (o ListRunsOptions) SetBefore(before string) ListRunsOptions {
	o["before"] = before
	return o
}

// ListRuns fetches runs with given `threadID` and `options`.
//
// https://platform.openai.com/docs/api-reference/runs/listRuns
func (c *Client) ListRuns(threadID string, options ListRunsOptions) (response Runs, err error) {
	if options == nil {
		options = ListRunsOptions{}
	}

	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/runs", threadID), options); err == nil {
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

	return Runs{}, err
}

// ToolOutput struct for API request
type ToolOutput struct {
	ToolCallID *string `json:"tool_call_id,omitempty"`
	Output     *string `json:"output,omitempty"`
}

// SubmitToolOutputs submits tool outputs with given `threadID` and `runID`.
//
// This can be called when
//
//	run.Status == RunStatusRequiresAction && run.RequiredAction.Type == "submit_tool_outputs".
//
// https://platform.openai.com/docs/api-reference/runs/submitToolOutputs
func (c *Client) SubmitToolOutputs(threadID, runID string, toolOutputs []ToolOutput) (response Run, err error) {
	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/threads/%s/runs/%s/submit_tool_outputs", threadID, runID), map[string]any{
		"tool_outputs": toolOutputs,
	}); err == nil {
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

	return Run{}, err
}

// CancelRun cancels a run with given `threadID` and `runID`.
//
// This can be called when
//
//	run.Status == RunStatusInProgress.
//
// https://platform.openai.com/docs/api-reference/runs/cancelRun
func (c *Client) CancelRun(threadID, runID string) (response Run, err error) {
	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/threads/%s/runs/%s/cancel", threadID, runID), nil); err == nil {
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

	return Run{}, err
}

// CreateThreadAndRunOptions for creating thread and running it
type CreateThreadAndRunOptions map[string]any

// SetThread sets the `thread` parameter of CreateThreadAndRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createThreadAndRun#runs-createthreadandrun-thread
func (o CreateThreadAndRunOptions) SetThread(thread RunnableThread) CreateThreadAndRunOptions {
	o["thread"] = thread
	return o
}

// SetModel sets the `model` parameter of CreateThreadAndRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createThreadAndRun#runs-createthreadandrun-model
func (o CreateThreadAndRunOptions) SetModel(model string) CreateThreadAndRunOptions {
	o["model"] = model
	return o
}

// SetInstructions sets the `instructions` parameter of CreateThreadAndRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createThreadAndRun#runs-createthreadandrun-instructions
func (o CreateThreadAndRunOptions) SetInstructions(instructions string) CreateThreadAndRunOptions {
	o["instructions"] = instructions
	return o
}

// SetTools sets the `tools` parameter of CreateThreadAndRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createThreadAndRun#runs-createthreadandrun-tools
func (o CreateThreadAndRunOptions) SetTools(tools []Tool) CreateThreadAndRunOptions {
	o["tools"] = tools
	return o
}

// SetMetadata sets the `metadata` parameter of CreateThreadAndRunOptions.
//
// https://platform.openai.com/docs/api-reference/runs/createThreadAndRun#runs-createthreadandrun-metadata
func (o CreateThreadAndRunOptions) SetMetadata(metadata map[string]string) CreateThreadAndRunOptions {
	o["metadata"] = metadata
	return o
}

// RunnableThread struct for CreateThreadAndRunOptions
type RunnableThread struct {
	Messages []RunnableThreadMessage `json:"messages,omitempty"`
	Metadata map[string]string       `json:"metadata,omitempty"`
}

// RunnableThreadMessage struct for RunnableThread struct
type RunnableThreadMessage struct {
	Role     string            `json:"role"` // == 'user'
	Content  string            `json:"content"`
	FileIDs  []string          `json:"file_ids,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CreateThreadAndRun creates a thread and runs it with given `assistantID` and `options`.
//
// https://platform.openai.com/docs/api-reference/runs/createThreadAndRun
func (c *Client) CreateThreadAndRun(assistantID string, options CreateThreadAndRunOptions) (response Run, err error) {
	if options == nil {
		options = CreateThreadAndRunOptions{}
	}
	options["assistant_id"] = assistantID

	var bytes []byte
	if bytes, err = c.post("v1/threads/runs", options); err == nil {
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

	return Run{}, err
}

// RunStepType type for constants
type RunStepType string

// RunStepType constants
const (
	RunStepTypeMessageCreation RunStepType = "message_creation"
	RunStepTypeToolCalls       RunStepType = "tool_calls"
)

// RunStepStatus type for constants
type RunStepStatus string

// RunStepStatus constants
const (
	RunStepStatusInProgress RunStepStatus = "in_progress"
	RunStepStatusCancelled  RunStepStatus = "cancelled"
	RunStepStatusFailed     RunStepStatus = "failed"
	RunStepStatusCompleted  RunStepStatus = "completed"
	RunStepStatusExpired    RunStepStatus = "expired"
)

// RunStepDetails struct for RunStepObject struct
type RunStepDetails struct {
	Type RunStepType `json:"type"`

	MessageCreation *RunStepDetailsMessageCreation `json:"message_creation,omitempty"` // Type == RunStepTypeMessageCreation
	ToolCalls       []RunStepDetailsToolCall       `json:"tool_calls,omitempty"`       // Type == RunStepTypeToolCalls
}

// RunStepDetailsMessageCreation struct for RunStepDetails struct
type RunStepDetailsMessageCreation struct {
	MessageID string `json:"message_id"`
}

// RunStepDetailsToolCall struct for RunStepDetails struct
type RunStepDetailsToolCall struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	CodeInterpreter *RunStepDetailsToolCallCodeInterpreter `json:"code_interpreter,omitempty"` // Type == ToolTypeCodeInterpreter
	Retrieval       *RunStepDetailsToolCallRetrieval       `json:"retrieval,omitempty"`        // Type -== ToolTypeRetrieval
	Function        *RunStepDetailsToolCallFunction        `json:"function,omitempty"`         // Type == ToolTypeFunction
}

// RunStepDetailsToolCallCodeInterpreter struct for RunStepDetailsToolCall struct
type RunStepDetailsToolCallCodeInterpreter struct {
	Input   string                          `json:"input"`
	Outputs []ToolCallCodeInterpreterOutput `json:"outputs"`
}

// ToolCallCodeInterpreterOutputType type for constants
type ToolCallCodeInterpreterOutputType string

// ToolCallCodeInterpreterOutputType constants
const (
	ToolCallCodeInterpreterOutputTypeLogs  ToolCallCodeInterpreterOutputType = "logs"
	ToolCallCodeInterpreterOutputTypeImage ToolCallCodeInterpreterOutputType = "image"
)

// ToolCallCodeInterpreterOutputImage struct for ToolCallCodeInterpreterOutput struct
type ToolCallCodeInterpreterOutputImage struct {
	FileID string `json:"file_id"`
}

// ToolCallCodeInterpreterOutput struct for RunStepDetailsToolCallCodeInterpreter struct
type ToolCallCodeInterpreterOutput struct {
	Type ToolCallCodeInterpreterOutputType `json:"type"`

	Logs  *string                             `json:"logs,omitempty"`  // Type == ToolCallCodeInterpreterOutputTypeLogs
	Image *ToolCallCodeInterpreterOutputImage `json:"image,omitempty"` // Type == ToolCallCodeInterpreterOutputTypeImage
}

// RunStepDetailsToolCallRetrieval struct for RunStepDetailsToolCall struct (empty object for now)
type RunStepDetailsToolCallRetrieval struct{}

// RunStepDetailsToolCallsFunction struct for RunStepDetailsToolCall struct
type RunStepDetailsToolCallFunction struct {
	Name      string  `json:"name"`
	Arguments string  `json:"arguments"`
	Output    *string `json:"output,omitempty"`
}

// https://platform.openai.com/docs/api-reference/runs/step-object
type RunStep struct {
	CommonResponse

	ID          string            `json:"id"`
	CreatedAt   int               `json:"created_at"`
	AssistantID string            `json:"assistant_id"`
	ThreadID    string            `json:"thread_id"`
	RunID       string            `json:"run_id"`
	Type        RunStepType       `json:"type"`
	Status      RunStepStatus     `json:"status"`
	StepDetails RunStepDetails    `json:"step_details"`
	LastError   *RunError         `json:"last_error,omitempty"`
	ExpiredAt   *int              `json:"expired_at,omitempty"`
	CancelledAt *int              `json:"cancelled_at,omitempty"`
	FailedAt    *int              `json:"failed_at,omitempty"`
	CompletedAt *int              `json:"completed_at,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

// RetrieveRunStep retrieves a run step with given `threadID`, `runID` and `stepID`.
//
// https://platform.openai.com/docs/api-reference/runs/getRunStep
func (c *Client) RetrieveRunStep(threadID, runID, stepID string) (response RunStep, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/runs/%s/steps/%s", threadID, runID, stepID), nil); err == nil {
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

	return RunStep{}, err
}

// ListRunStepsOptions type for listing run steps
type ListRunStepsOptions map[string]any

// SetLimit sets the `limit` parameter of run steps' listing request.
//
// https://platform.openai.com/docs/api-reference/runs/listRunSteps#runs-listrunsteps-limit
func (o ListRunStepsOptions) SetLimit(limit int) ListRunStepsOptions {
	o["limit"] = limit
	return o
}

// SetOrder sets the `order` parameter of run steps' listing request.
//
// `order` can be one of 'asc' or 'desc'. (default: 'desc')
//
// https://platform.openai.com/docs/api-reference/runs/listRunSteps#runs-listrunsteps-order
func (o ListRunStepsOptions) SetOrder(order string) ListRunStepsOptions {
	o["order"] = order
	return o
}

// SetAfter sets the `after` parameter of run steps' listing request.
//
// https://platform.openai.com/docs/api-reference/runs/listRunSteps#runs-listrunsteps-after
func (o ListRunStepsOptions) SetAfter(after string) ListRunStepsOptions {
	o["after"] = after
	return o
}

// SetBefore sets the `before` parameter of run steps' listing request.
//
// https://platform.openai.com/docs/api-reference/runs/listRunSteps#runs-listrunsteps-before
func (o ListRunStepsOptions) SetBefore(before string) ListRunStepsOptions {
	o["before"] = before
	return o
}

// RunSteps struct for API response
type RunSteps struct {
	CommonResponse

	Data    []RunStep `json:"data"`
	FirstID string    `json:"first_id"`
	LastID  string    `json:"last_id"`
	HasMore bool      `json:"has_more"`
}

// ListRunSteps fetches run steps with given `threadID`, `runID` and `options`.
//
// https://platform.openai.com/docs/api-reference/runs/listRunSteps
func (c *Client) ListRunSteps(threadID, runID string, options ListRunStepsOptions) (response RunSteps, err error) {
	if options == nil {
		options = ListRunStepsOptions{}
	}

	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/runs/%s/steps", threadID, runID), options); err == nil {
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

	return RunSteps{}, err
}
