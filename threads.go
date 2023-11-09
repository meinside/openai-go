package openai

import (
	"encoding/json"
	"fmt"
)

// NOTE: In Beta

// https://platform.openai.com/docs/api-reference/threads

// Thread struct
//
// https://platform.openai.com/docs/api-reference/threads/object
type Thread struct {
	CommonResponse

	ID        string            `json:"id"`
	CreatedAt int               `json:"created_at"`
	Metadata  map[string]string `json:"metadata"`
}

// ThreadMessage struct for Thread
type ThreadMessage struct {
	Role     string            `json:"role"` // == 'user'
	Content  string            `json:"content"`
	FileIDs  []string          `json:"file_ids,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// NewThreadMessage returns a new ThreadMessage with given `content`.
func NewThreadMessage(content string) ThreadMessage {
	return ThreadMessage{
		Role:    "user",
		Content: content,
	}
}

// SetFileIDs sets the `file_ids` value of ThreadMessage and return it.
func (m ThreadMessage) SetFileIDs(fileIDs []string) ThreadMessage {
	m.FileIDs = fileIDs
	return m
}

// SetMetadata sets the `metadata` value of ThreadMessage and return it.
func (m ThreadMessage) SetMetadata(metadata map[string]string) ThreadMessage {
	m.Metadata = metadata
	return m
}

// CreateThreadOptions for creating thread
type CreateThreadOptions map[string]any

// SetMessages sets the `messages` parameter of CreateThreadOptions.
//
// https://platform.openai.com/docs/api-reference/threads/createThread#threads-createthread-messages
func (o CreateThreadOptions) SetMessages(messages []ThreadMessage) CreateThreadOptions {
	o["messages"] = messages
	return o
}

// SetMetadata sets the `metadata` parameter of CreateThreadOptions.
//
// https://platform.openai.com/docs/api-reference/threads/createThread#threads-createthread-metadata
func (o CreateThreadOptions) SetMetadata(metadata map[string]string) CreateThreadOptions {
	o["metadata"] = metadata
	return o
}

// CreateThread creates a thread with given `options`.
//
// https://platform.openai.com/docs/api-reference/threads/createThread
func (c *Client) CreateThread(options CreateThreadOptions) (response Thread, err error) {
	if options == nil {
		options = CreateThreadOptions{}
	}

	var bytes []byte
	if bytes, err = c.post("v1/threads", options); err == nil {
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

	return Thread{}, err
}

// RetrieveThread retrieves the thread with given `threadID`.
//
// https://platform.openai.com/docs/api-reference/threads/getThread
func (c *Client) RetrieveThread(threadID string) (response Thread, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s", threadID), nil); err == nil {
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

	return Thread{}, err
}

// ModifyThreadOptions for modifying thread
type ModifyThreadOptions map[string]any

// SetMetadata sets the `metadata` parameter of ModifyThreadOptions.
//
// https://platform.openai.com/docs/api-reference/threads/modifyThread#threads-modifythread-metadata
func (o ModifyThreadOptions) SetMetadata(metadata map[string]string) ModifyThreadOptions {
	o["metadata"] = metadata
	return o
}

// ModifyThread modifies a thread with given `threadID` and `options`.
//
// https://platform.openai.com/docs/api-reference/threads/modifyThread
func (c *Client) ModifyThread(threadID string, options ModifyThreadOptions) (response Thread, err error) {
	if options == nil {
		options = ModifyThreadOptions{}
	}

	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/threads/%s", threadID), options); err == nil {
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

	return Thread{}, err
}

// ThreadDeletionStatus for API response
type ThreadDeletionStatus struct {
	CommonResponse

	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// DeleteThread deletes a thread with given `threadID`.
//
// https://platform.openai.com/docs/api-reference/threads/deleteThread
func (c *Client) DeleteThread(threadID string) (response ThreadDeletionStatus, err error) {
	var bytes []byte
	if bytes, err = c.delete(fmt.Sprintf("v1/threads/%s", threadID), nil); err == nil {
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

	return ThreadDeletionStatus{}, err
}
