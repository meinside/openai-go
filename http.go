package openai

// types and functions for HTTP requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"os"
	"strings"
)

const (
	baseURL            = "https://api.openai.com"
	defaultContentType = "application/json"

	kContentType        = "Content-Type"
	kContentDisposition = "Content-Disposition"
	kAuthorization      = "Authorization"
	kOrganization       = "OpenAI-Organization"
)

// CommonResponse struct for responses with common properties
type CommonResponse struct {
	Object *string `json:"object,omitempty"`
	Error  *Error  `json:"error,omitempty"`
}

// Error struct for response error property
type Error struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   any     `json:"param,omitempty"`
	Code    *string `json:"code,omitempty"`
}

// Usage struct for reponses
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// err converts `Error` to `error`.
func (e *Error) err() error {
	es := map[string]any{
		"type":    e.Type,
		"message": e.Message,
	}
	if e.Code != nil {
		es["code"] = *e.Code
	}
	if e.Param != nil {
		es["param"] = e.Param
	}

	if bytes, err := json.Marshal(es); err == nil {
		return fmt.Errorf(string(bytes))
	} else {
		return fmt.Errorf(fmt.Sprintf("%+v", es))
	}
}

// FileParam struct for multipart requests
type FileParam struct {
	bs []byte
}

// NewFileParamFromBytes returns a new FileParam with given bytes
func NewFileParamFromBytes(bs []byte) FileParam {
	return FileParam{
		bs: bs,
	}
}

// NewFileParamFromFilepath returns a new FileParam with bytes read from given filepath
func NewFileParamFromFilepath(path string) (f FileParam, err error) {
	var bs []byte
	if bs, err = os.ReadFile(path); err == nil {
		return FileParam{
			bs: bs,
		}, nil
	}
	return FileParam{}, err
}

// sends HTTP request
func (c *Client) do(method, endpoint string, params map[string]any) (response []byte, err error) {
	if params == nil {
		params = map[string]any{}
	}

	apiURL := fmt.Sprintf("%s/%s", baseURL, endpoint)

	var req *http.Request
	if req, err = http.NewRequest(method, apiURL, nil); err == nil {
		// parameters
		queries := req.URL.Query()
		for k, v := range params {
			queries.Add(k, fmt.Sprintf("%+v", v))
		}
		req.URL.RawQuery = queries.Encode()

		// headers
		req.Header.Set(kAuthorization, fmt.Sprintf("Bearer %s", c.APIKey))
		req.Header.Set(kOrganization, c.OrganizationID)

		if c.Verbose {
			if dumped, err := httputil.DumpRequest(req, true); err == nil {
				log.Printf("dump request:\n\n%s", string(dumped))
			}
		}

		req.Close = true

		// send request and return response bytes
		var resp *http.Response
		resp, err = c.httpClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err == nil {
			if response, err = io.ReadAll(resp.Body); err == nil {
				if c.Verbose {
					log.Printf("API response for %s: '%s'", endpoint, string(response))
				}

				if resp.StatusCode != 200 {
					err = fmt.Errorf("http status %d", resp.StatusCode)
				}

				return response, err
			}
		}
	}

	return nil, err
}

// sends HTTP GET request
func (c *Client) get(endpoint string, params map[string]any) (response []byte, err error) {
	return c.do(http.MethodGet, endpoint, params)
}

// sends HTTP DELETE request
func (c *Client) delete(endpoint string, params map[string]any) (response []byte, err error) {
	return c.do(http.MethodDelete, endpoint, params)
}

// sends HTTP POST request
func (c *Client) post(endpoint string, params map[string]any) (response []byte, err error) {
	if params == nil {
		params = map[string]any{}
	}

	apiURL := fmt.Sprintf("%s/%s", baseURL, endpoint)

	var req *http.Request

	if hasFileInParams(params) {
		// multipart/form-data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for k, v := range params {
			switch val := v.(type) {
			case FileParam:
				bs := val.bs
				filename := fmt.Sprintf("%s.%s", k, getExtension(bs))

				var part io.Writer
				if part, err = writer.CreatePart(mimeHeaderForBytes(bs, k, filename)); err == nil {
					if _, err = io.Copy(part, bytes.NewReader(bs)); err != nil {
						return nil, fmt.Errorf("could not write bytes to multipart for param '%s': %s", k, err)
					}
				} else {
					return nil, fmt.Errorf("could not create part for param '%s': %s", k, err)
				}
			default:
				if err := writer.WriteField(k, fmt.Sprintf("%v", v)); err != nil {
					return nil, fmt.Errorf("could not write field with key: %s, value: %v", k, v)
				}
			}
		}

		if err = writer.Close(); err != nil {
			return nil, fmt.Errorf("error while closing multipart form data writer: %s", err)
		}

		if req, err = http.NewRequest(http.MethodPost, apiURL, body); err != nil {
			return nil, fmt.Errorf("failed to create multipart request: %s", err)
		}

		// set content-type header
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		// application/json
		var serialized []byte
		if serialized, err = json.Marshal(params); err == nil {
			if req, err = http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(serialized)); err != nil {
				return nil, fmt.Errorf("failed to create application/json request: %s", err)
			}

			// set content-type header
			req.Header.Set(kContentType, defaultContentType)
		}
	}

	// set authentication headers
	req.Header.Set(kAuthorization, fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set(kOrganization, c.OrganizationID)

	if c.Verbose {
		if dumped, err := httputil.DumpRequest(req, true); err == nil {
			log.Printf("dump request:\n\n%s", string(dumped))
		}
	}

	req.Close = true

	// send request and return response bytes
	var resp *http.Response
	resp, err = c.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err == nil {
		if response, err = io.ReadAll(resp.Body); err == nil {
			if c.Verbose {
				log.Printf("API response for %s: '%s'", endpoint, string(response))
			}

			if resp.StatusCode != 200 {
				err = fmt.Errorf("http status %d", resp.StatusCode)
			}

			return response, err
		}
	}

	return nil, err
}

// checks if given params include any file param
func hasFileInParams(params map[string]any) bool {
	for _, v := range params {
		if _, ok := v.(FileParam); ok {
			return true
		}
	}
	return false
}

// get file extension from given bytes array
//
// https://www.w3.org/Protocols/rfc1341/4_Content-Type.html
func getExtension(bytes []byte) string {
	types := strings.Split(http.DetectContentType(bytes), "/") // ex: "image/jpeg"
	if len(types) >= 2 {
		splitted := strings.Split(types[1], ";") // for removing subtype parameter
		if len(splitted) >= 1 {
			if splitted[0] == "wave" {
				return "wav"
			}

			return splitted[0] // return subtype only
		}
	}
	return ""
}

// generates mime header
func mimeHeaderForBytes(bs []byte, key, filename string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set(kContentDisposition, fmt.Sprintf(`form-data; name="%s"; filename="%s"`, key, filename))
	h.Set(kContentType, http.DetectContentType(bs))
	return h
}
