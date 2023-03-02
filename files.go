package openai

import (
	"encoding/json"
	"fmt"
)

// https://platform.openai.com/docs/api-reference/files

// File struct for file items
type File struct {
	CommonResponse

	ID        string `json:"id"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// UploadedFile struct for response
type UploadedFile File

// RetrievedFile struct for response
type RetrievedFile File

// DeletedFile struct for response
type DeletedFile struct {
	CommonResponse

	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// Files struct for response
type Files struct {
	CommonResponse

	Data []File `json:"data"`
}

// ListFiles returns a list of files that belong to the requested organization id.
//
// https://platform.openai.com/docs/api-reference/files/list
func (c *Client) ListFiles() (response Files, err error) {
	var bytes []byte
	if bytes, err = c.get("v1/files", nil); err == nil {
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

	return Files{}, err
}

// UploadFile uploads given file.
//
// https://platform.openai.com/docs/api-reference/files/upload
func (c *Client) UploadFile(file FileParam, purpose string) (response UploadedFile, err error) {
	var bytes []byte
	if bytes, err = c.post("v1/files", map[string]any{
		"file":    file,
		"purpose": purpose,
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

	return UploadedFile{}, err
}

// DeleteFile deletes given file.
//
// https://platform.openai.com/docs/api-reference/files/delete
func (c *Client) DeleteFile(fileID string) (response DeletedFile, err error) {
	var bytes []byte
	if bytes, err = c.delete(fmt.Sprintf("v1/files/%s", fileID), nil); err == nil {
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

	return DeletedFile{}, err
}

// RetrieveFile returns the information of given file.
//
// https://platform.openai.com/docs/api-reference/files/retrieve
func (c *Client) RetrieveFile(fileID string) (response RetrievedFile, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/files/%s", fileID), nil); err == nil {
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

	return RetrievedFile{}, err
}

// RetrieveFileContent returns the content of given file.
//
// https://platform.openai.com/docs/api-reference/files/retrieve-content
func (c *Client) RetrieveFileContent(fileID string) (response []byte, err error) {
	var bytes []byte
	if bytes, err = c.get(fmt.Sprintf("v1/files/%s/content", fileID), nil); err == nil {
		return bytes, nil
	} else {
		var res CommonResponse
		if e := json.Unmarshal(bytes, &res); e == nil {
			err = fmt.Errorf("%s: %s", err, res.Error.err())
		}
	}

	return nil, err
}
