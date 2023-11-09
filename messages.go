package openai

import (
	"encoding/json"
	"fmt"
)

// NOTE: In Beta

// https://platform.openai.com/docs/api-reference/messages

// https://platform.openai.com/docs/api-reference/messages/object
type Message struct {
	CommonResponse

	ID          string            `json:"id"`
	CreatedAt   int               `json:"created_at"`
	ThreadID    string            `json:"thread_id"`
	Role        string            `json:"role"` // 'user' | 'assistant'
	Content     []MessageContent  `json:"content"`
	AssistantID *string           `json:"assistant_id,omitempty"`
	RunID       *string           `json:"run_id,omitempty"`
	FileIDs     []string          `json:"file_ids"`
	Metadata    map[string]string `json:"metadata"`
}

// MessageContentType type for constants
type MessageContentType string

// MessageContentType constants
const (
	MessageContentTypeImageFile MessageContentType = "image_file"
	MessageContentTypeText      MessageContentType = "text"
)

// MessageContent struct for Message
type MessageContent struct {
	Type MessageContentType `json:"type"`

	ImageFile *MessageContentImageFile `json:"image_file,omitempty"` // Type == 'image_file'
	Text      *MessageContentText      `json:"text,omitempty"`       // Type == 'text'
}

// MessageContentImageFile struct for MessageContent
type MessageContentImageFile struct {
	FileID string `json:"file_id"`
}

// MessageContentText struct for MessageContent
type MessageContentText struct {
	Value       string                         `json:"value"`
	Annotations []MessageContentTextAnnotation `json:"annotations"`
}

// MessageContentTextAnnotationType type for constants
type MessageContentTextAnnotationType string

// MessageContentTextAnnotationType constants
const (
	MessageContentTextAnnotationTypeFileCitation MessageContentTextAnnotationType = "file_citation"
	MessageContentTextAnnotationTypeFilePath     MessageContentTextAnnotationType = "file_path"
)

// MessageContentTextAnntation struct for MessageContentText
type MessageContentTextAnnotation struct {
	Type MessageContentTextAnnotationType `json:"type"`
	Text string                           `json:"text"`

	FileCitation *MessageContentTextAnnotationFileCitation `json:"file_citation,omitempty"`
	FilePath     *MessageContentTextAnnotationFilePath     `json:"file_path,omitempty"`

	StartIndex int `json:"start_index"`
	EndIndex   int `json:"end_index"`
}

// MessageContentTextAnnotationFileCitation struct
type MessageContentTextAnnotationFileCitation struct {
	FileID string `json:"file_id"`
	Quote  string `json:"quote"`
}

// MessageContentTextAnnotationFilePath struct
type MessageContentTextAnnotationFilePath struct {
	FileID string `json:"file_id"`
}

// CreateMessageOptions for creating message
type CreateMessageOptions map[string]any

// SetFileIDs sets the `file_ids` parameter of CreateMessageOptions.
//
// https://platform.openai.com/docs/api-reference/messages/createMessage#messages-createmessage-file_ids
func (o CreateMessageOptions) SetFileIDs(fileIDs []string) CreateMessageOptions {
	o["file_ids"] = fileIDs
	return o
}

// SetMetadata sets the `metadata` parameter of CreateMessageOptions.
//
// https://platform.openai.com/docs/api-reference/messages/createMessage#messages-createmessage-metadata
func (o CreateMessageOptions) SetMetadata(metadata map[string]string) CreateMessageOptions {
	o["metadata"] = metadata
	return o
}

// CreateMessage creates a message with given `threadID`, `role`, `content`, and `options`.
//
// https://platform.openai.com/docs/api-reference/messages/createMessage
func (c *Client) CreateMessage(threadID, role, content string, options CreateMessageOptions) (response Message, err error) {
	if options == nil {
		options = CreateMessageOptions{}
	}
	options["role"] = role
	options["content"] = content

	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/threads/%s/messages", threadID), options); err == nil {
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

	return Message{}, err
}

// RetrieveMessage retrieves a message with given `threadID` and `messageID`.
//
// https://platform.openai.com/docs/api-reference/messages/getMessage
func (c *Client) RetrieveMessage(threadID, messageID string) (response Message, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/messages/%s", threadID, messageID), nil); err == nil {
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

	return Message{}, err
}

// ModifyMessageOptions for modifying message
type ModifyMessageOptions map[string]any

// SetMetadata sets the `metadata` parameter of ModifyMessageOptions.
//
// https://platform.openai.com/docs/api-reference/messages/modifyMessage#messages-modifymessage-metadata
func (o ModifyMessageOptions) SetMetadata(metadata map[string]string) ModifyMessageOptions {
	o["metadata"] = metadata
	return o
}

// ModifyMessage modifies a message with given `threadID`, `messageID`, and `options`.
//
// https://platform.openai.com/docs/api-reference/messages/modifyMessage
func (c *Client) ModifyMessage(threadID, messageID string, options ModifyMessageOptions) (response Message, err error) {
	if options == nil {
		options = ModifyMessageOptions{}
	}

	var bytes []byte
	if bytes, err = c.post(fmt.Sprintf("v1/threads/%s/messages/%s", threadID, messageID), options); err == nil {
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

	return Message{}, err
}

// Messages struct for API response
type Messages struct {
	CommonResponse

	Data    []Message `json:"data"`
	FirstID string    `json:"first_id"`
	LastID  string    `json:"last_id"`
	HasMore bool      `json:"has_more"`
}

// ListMessagesOptions for listing messages
type ListMessagesOptions map[string]any

// SetLimit sets the `limit` parameter of messages' listing request.
//
// https://platform.openai.com/docs/api-reference/messages/listMessages#messages-listmessages-limit
func (o ListMessagesOptions) SetLimit(limit int) ListMessagesOptions {
	o["limit"] = limit
	return o
}

// SetOrder sets the `order` parameter of messages' listing request.
//
// `order` can be one of 'asc' or 'desc'. (default: 'desc')
//
// https://platform.openai.com/docs/api-reference/messages/listMessages#messages-listmessages-order
func (o ListMessagesOptions) SetOrder(order string) ListMessagesOptions {
	o["order"] = order
	return o
}

// SetAfter sets the `after` parameter of messages' listing request.
//
// https://platform.openai.com/docs/api-reference/messages/listMessages#messages-listmessages-after
func (o ListMessagesOptions) SetAfter(after string) ListMessagesOptions {
	o["after"] = after
	return o
}

// SetBefore sets the `before` parameter of messages' listing request.
//
// https://platform.openai.com/docs/api-reference/messages/listMessages#messages-listmessages-before
func (o ListMessagesOptions) SetBefore(before string) ListMessagesOptions {
	o["before"] = before
	return o
}

// ListMessages fetches messages with given `threadID`, and `options`.
//
// https://platform.openai.com/docs/api-reference/messages/listMessages
func (c *Client) ListMessages(threadID string, options ListMessagesOptions) (response Messages, err error) {
	if options == nil {
		options = ListMessagesOptions{}
	}

	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/messages", threadID), options); err == nil {
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

	return Messages{}, err
}

// https://platform.openai.com/docs/api-reference/messages/file-object
type MessageFile struct {
	CommonResponse

	ID        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	MessageID string `json:"message_id"`
}

// RetrieveMessageFile retrieves a message file with given `threadID`, `messageID`, and `fileID`.
//
// https://platform.openai.com/docs/api-reference/messages/getMessageFile
func (c *Client) RetrieveMessageFile(threadID, messageID, fileID string) (response MessageFile, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/messages/%s/files/%s", threadID, messageID, fileID), nil); err == nil {
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

	return MessageFile{}, err
}

// MessageFiles struct for API response
type MessageFiles struct {
	CommonResponse

	Data    []MessageFile `json:"data"`
	FirstID string        `json:"first_id"`
	LastID  string        `json:"last_id"`
	HasMore bool          `json:"has_more"`
}

// ListMessageFilesOptions for listing message files
type ListMessageFilesOptions map[string]any

// SetLimit sets the `limit` parameter of message files' listing request.
//
// https://platform.openai.com/docs/api-reference/messages/listMessageFiles#messages-listmessagefiles-limit
func (o ListMessageFilesOptions) SetLimit(limit int) ListMessageFilesOptions {
	o["limit"] = limit
	return o
}

// SetOrder sets the `order` parameter of message files' listing request.
//
// `order` can be one of 'asc' or 'desc'. (default: 'desc')
//
// https://platform.openai.com/docs/api-reference/messages/listMessageFiles#messages-listmessagefiles-order
func (o ListMessageFilesOptions) SetOrder(order string) ListMessageFilesOptions {
	o["order"] = order
	return o
}

// SetAfter sets the `after` parameter of message files' listing request.
//
// https://platform.openai.com/docs/api-reference/messages/listMessageFiles#messages-listmessagefiles-after
func (o ListMessageFilesOptions) SetAfter(after string) ListMessageFilesOptions {
	o["after"] = after
	return o
}

// SetBefore sets the `before` parameter of messages' listing request.
//
// https://platform.openai.com/docs/api-reference/messages/listMessageFiles#messages-listmessagefiles-before
func (o ListMessageFilesOptions) SetBefore(before string) ListMessageFilesOptions {
	o["before"] = before
	return o
}

// ListMessageFiles fetches message files with given `threadID`, `mesageID`, and `options`.
//
// https://platform.openai.com/docs/api-reference/messages/listMessageFiles
func (c *Client) ListMessageFiles(threadID, messageID string, options ListMessageFilesOptions) (response MessageFiles, err error) {
	if options == nil {
		options = ListMessageFilesOptions{}
	}

	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/threads/%s/messages/%s/files", threadID, messageID), options); err == nil {
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

	return MessageFiles{}, err
}
