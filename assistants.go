package openai

import (
	"encoding/json"
	"fmt"
)

// NOTE: In Beta

// https://platform.openai.com/docs/api-reference/assistants

// Assistant struct for assistant object
//
// https://platform.openai.com/docs/api-reference/assistants/object
type Assistant struct {
	CommonResponse

	ID           string            `json:"id"`
	CreatedAt    int               `json:"created_at"`
	Name         *string           `json:"name,omitempty"`
	Description  *string           `json:"description,omitempty"`
	Model        string            `json:"model"`
	Instructions *string           `json:"instructions,omitempty"`
	Tools        []Tool            `json:"tools"`
	FileIDs      []string          `json:"file_ids"`
	Metadata     map[string]string `json:"metadata"`
}

// Tool struct for assistant object
//
// https://platform.openai.com/docs/api-reference/assistants/object#assistants/object-tools
type Tool struct {
	Type     string        `json:"type"`
	Function *ToolFunction `json:"function,omitempty"`
}

// NewCodeInterpreterTool returns a tool with type: 'code_interpreter'.
func NewCodeInterpreterTool() Tool {
	return Tool{
		Type: "code_interpreter",
	}
}

// NewRetrievalTool returns a tool with type: 'retrieval'.
func NewRetrievalTool() Tool {
	return Tool{
		Type: "retrieval",
	}
}

// NewFunctionTool returns a tool with type: 'function'.
func NewFunctionTool(fn ToolFunction) Tool {
	return Tool{
		Type:     "function",
		Function: &fn,
	}
}

// NewBuiltinTool returns a tool with a specified type.
func NewBuiltinTool(fn string) Tool {
	return Tool{
		Type: fn,
	}
}

// ToolFunction struct for Tool struct
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  ToolFunctionParameters `json:"parameters"`
}

// CreateAssistantOptions for creating assistant
type CreateAssistantOptions map[string]any

// SetName sets the `name` parameter of assistant creation.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistant#assistants-createassistant-name
func (o CreateAssistantOptions) SetName(name string) CreateAssistantOptions {
	o["name"] = name
	return o
}

// SetDescription sets the `description` parameter of assistant creation.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistant#assistants-createassistant-description
func (o CreateAssistantOptions) SetDescription(description string) CreateAssistantOptions {
	o["description"] = description
	return o
}

// SetInstructions sets the `instructions` parameter of assistant creation.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistant#assistants-createassistant-instructions
func (o CreateAssistantOptions) SetInstructions(instructions string) CreateAssistantOptions {
	o["instructions"] = instructions
	return o
}

// SetTools sets the `tools` parameter of assistant creation.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistant#assistants-createassistant-tools
func (o CreateAssistantOptions) SetTools(tools []Tool) CreateAssistantOptions {
	o["tools"] = tools
	return o
}

// SetFileIDs sets the `file_ids` parameter of assistant creation.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistant#assistants-createassistant-file_ids
func (o CreateAssistantOptions) SetFileIDs(fileIDs []string) CreateAssistantOptions {
	o["file_ids"] = fileIDs
	return o
}

// SetMetadata sets the `metadata` parameter of assistant creation.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistant#assistants-createassistant-metadata
func (o CreateAssistantOptions) SetMetadata(metadata map[string]string) CreateAssistantOptions {
	o["metadata"] = metadata
	return o
}

// CreateAssistant creates an assitant with given `model` and `options`.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistant
func (c *Client) CreateAssistant(model string, options CreateAssistantOptions) (response Assistant, err error) {
	if options == nil {
		options = CreateAssistantOptions{}
	}
	options["model"] = model

	var bytes []byte
	if bytes, err = c.post("v1/assistants", options); err == nil {
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

	return Assistant{}, err
}

// RetrieveAssistant retrieves an assistant with given `assistantID`.
//
// https://platform.openai.com/docs/api-reference/assistants/getAssistant
func (c *Client) RetrieveAssistant(assistantID string) (response Assistant, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/assistants/%s", assistantID), nil); err == nil {
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

	return Assistant{}, err
}

// ModifyAssistantOptions for modifying assistant
type ModifyAssistantOptions map[string]any

// SetModel sets the `model` parameter of assistant modification.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant#assistants-modifyassistant-model
func (o ModifyAssistantOptions) SetModel(model string) ModifyAssistantOptions {
	o["model"] = model
	return o
}

// SetName sets the `name` parameter of assistant modification.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant#assistants-modifyassistant-name
func (o ModifyAssistantOptions) SetName(name string) ModifyAssistantOptions {
	o["name"] = name
	return o
}

// SetDescription sets the `description` parameter of assistant modification.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant#assistants-modifyassistant-description
func (o ModifyAssistantOptions) SetDescription(description string) ModifyAssistantOptions {
	o["description"] = description
	return o
}

// SetInstructions sets the `instructions` parameter of assistant modification.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant#assistants-modifyassistant-instructions
func (o ModifyAssistantOptions) SetInstructions(instructions string) ModifyAssistantOptions {
	o["instructions"] = instructions
	return o
}

// SetTools sets the `tools` parameter of assistant modification.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant#assistants-modifyassistant-tools
func (o ModifyAssistantOptions) SetTools(tools []Tool) ModifyAssistantOptions {
	o["tools"] = tools
	return o
}

// SetFileIDs sets the `file_ids` parameter of assistant modification.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant#assistants-modifyassistant-file_ids
func (o ModifyAssistantOptions) SetFileIDs(fileIDs []string) ModifyAssistantOptions {
	o["file_ids"] = fileIDs
	return o
}

// SetMetadata sets the `metadata` parameter of assistant modification.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant#assistants-modifyassistant-metadata
func (o ModifyAssistantOptions) SetMetadata(metadata map[string]string) ModifyAssistantOptions {
	o["metadata"] = metadata
	return o
}

// ModifyAssistant modifies an assistant with given `assistantID` and `options`.
//
// https://platform.openai.com/docs/api-reference/assistants/modifyAssistant
func (c *Client) ModifyAssistant(assistantID string, options ModifyAssistantOptions) (response Assistant, err error) {
	if options == nil {
		options = ModifyAssistantOptions{}
	}

	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/assistants/%s", assistantID), options); err == nil {
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

	return Assistant{}, err
}

// AssistantDeletionStatus struct for API response
type AssistantDeletionStatus struct {
	CommonResponse

	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// DeleteAssistant deletes an assistant with given `assistantID`.
//
// https://platform.openai.com/docs/api-reference/assistants/deleteAssistant
func (c *Client) DeleteAssistant(assistantID string) (response AssistantDeletionStatus, err error) {
	var bytes []byte
	if bytes, err = c.delete(fmt.Sprintf("v1/assistants/%s", assistantID), nil); err == nil {
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

	return AssistantDeletionStatus{}, err
}

// ListAssistantsOptions for listing assistants
type ListAssistantsOptions map[string]any

// SetLimit sets the `limit` parameter of assistants' listing request.
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistants#assistants-listassistants-limit
func (o ListAssistantsOptions) SetLimit(limit int) ListAssistantsOptions {
	o["limit"] = limit
	return o
}

// SetOrder sets the `order` parameter of assistants' listing request.
//
// `order` can be one of 'asc' or 'desc'. (default: 'desc')
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistants#assistants-listassistants-order
func (o ListAssistantsOptions) SetOrder(order string) ListAssistantsOptions {
	o["order"] = order
	return o
}

// SetAfter sets the `after` parameter of assistants' listing request.
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistants#assistants-listassistants-after
func (o ListAssistantsOptions) SetAfter(after string) ListAssistantsOptions {
	o["after"] = after
	return o
}

// SetBefore sets the `before` parameter of assistants' listing request.
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistants#assistants-listassistants-before
func (o ListAssistantsOptions) SetBefore(before string) ListAssistantsOptions {
	o["before"] = before
	return o
}

// Assistants struct for API response
type Assistants struct {
	CommonResponse

	Data    []Assistant `json:"data"`
	FirstID string      `json:"first_id"`
	LastID  string      `json:"last_id"`
	HasMore bool        `json:"has_more"`
}

// ListAssistants lists all assistants with given `options`.
//
// https://platform.openai.com/docs/api-reference/assistants/getAssistants
func (c *Client) ListAssistants(options ListAssistantsOptions) (response Assistants, err error) {
	if options == nil {
		options = ListAssistantsOptions{}
	}

	var bytes []byte
	if bytes, err = c.get("v1/assistants", options); err == nil {
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

	return Assistants{}, err
}

// AssistantFile struct for attached files of assistants
//
// https://platform.openai.com/docs/api-reference/assistants/file-object
type AssistantFile struct {
	CommonResponse

	ID          string `json:"id"`
	CreatedAt   int    `json:"created_at"`
	AssistantID string `json:"assistant_id"`
}

// CreateAssistantFile creates an assistant file by attaching given `fileID` to an assistant with `assistantID`.
//
// https://platform.openai.com/docs/api-reference/assistants/createAssistantFile
func (c *Client) CreateAssistantFile(assistantID, fileID string) (response AssistantFile, err error) {
	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/assistants/%s/files", assistantID), map[string]any{
		"file_id": fileID,
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

	return AssistantFile{}, err
}

// RetrieveAssistantFile retrieves an assistant file by given `assistantID` and `fileID`.
//
// https://platform.openai.com/docs/api-reference/assistants/getAssistantFile
func (c *Client) RetrieveAssistantFile(assistantID, fileID string) (response AssistantFile, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/assistants/%s/files/%s", assistantID, fileID), nil); err == nil {
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

	return AssistantFile{}, err
}

// AssistantFileDeletionStatus struct for API response
type AssistantFileDeletionStatus struct {
	CommonResponse

	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// DeleteAssistantFile deletes an assistant file by given `assistantID` and `fileID`.
//
// https://platform.openai.com/docs/api-reference/assistants/deleteAssistantFile
func (c *Client) DeleteAssistantFile(assistantID, fileID string) (response AssistantFileDeletionStatus, err error) {
	var bytes []byte
	if bytes, err = c.delete(fmt.Sprintf("v1/assistants/%s/files/%s", assistantID, fileID), nil); err == nil {
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

	return AssistantFileDeletionStatus{}, err
}

// AssistantFiless struct for API response
type AssistantFiles struct {
	CommonResponse

	Data    []AssistantFile `json:"data"`
	FirstID string          `json:"first_id"`
	LastID  string          `json:"last_id"`
	HasMore bool            `json:"has_more"`
}

// ListAssistantFilesOptions for listing assistant files
type ListAssistantFilesOptions map[string]any

// SetLimit sets the `limit` parameter of assistant files' listing request.
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistantFiles#assistants-listassistantfiles-limit
func (o ListAssistantFilesOptions) SetLimit(limit int) ListAssistantFilesOptions {
	o["limit"] = limit
	return o
}

// SetOrder sets the `order` parameter of assistant files' listing request.
//
// `order` can be one of 'asc' or 'desc'. (default: 'desc')
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistantFiles#assistants-listassistantfiles-order
func (o ListAssistantFilesOptions) SetOrder(order string) ListAssistantFilesOptions {
	o["order"] = order
	return o
}

// SetAfter sets the `after` parameter of assistant files' listing request.
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistantFiles#assistants-listassistantfiles-after
func (o ListAssistantFilesOptions) SetAfter(after string) ListAssistantFilesOptions {
	o["after"] = after
	return o
}

// SetBefore sets the `before` parameter of assistant files' listing request.
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistantFiles#assistants-listassistantfiles-before
func (o ListAssistantFilesOptions) SetBefore(before string) ListAssistantFilesOptions {
	o["before"] = before
	return o
}

// ListAssistantFiles lists all assistant files with given `assistantID` and `options`.
//
// https://platform.openai.com/docs/api-reference/assistants/listAssistantFiles
func (c *Client) ListAssistantFiles(assistantID string, options ListAssistantFilesOptions) (response AssistantFiles, err error) {
	if options == nil {
		options = ListAssistantFilesOptions{}
	}

	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/assistants/%s/files", assistantID), options); err == nil {
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

	return AssistantFiles{}, err
}
