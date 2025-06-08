package openai

// https://platform.openai.com/docs/api-reference/responses

import (
	"context"
	"encoding/json"
	"fmt"
)

// Response represents a response from the OpenAI Responses API
type Response struct {
	CommonResponse

	ID                 string           `json:"id"`
	Object             string           `json:"object"`
	CreatedAt          int64            `json:"created_at"`
	Status             string           `json:"status"`
	Error              *Error           `json:"error"`
	IncompleteDetails  any              `json:"incomplete_details"`
	Instructions       string           `json:"instructions,omitempty"`
	MaxOutputTokens    *int             `json:"max_output_tokens"`
	Model              string           `json:"model"`
	Output             []ResponseOutput `json:"output"`
	ParallelToolCalls  *bool            `json:"parallel_tool_calls,omitempty"`
	PreviousResponseID *string          `json:"previous_response_id"`
	Reasoning          any              `json:"reasoning,omitempty"`
	Store              *bool            `json:"store,omitempty"`
	Temperature        *float64         `json:"temperature,omitempty"`
	Text               any              `json:"text,omitempty"`
	ToolChoice         any              `json:"tool_choice,omitempty"`
	Tools              []any            `json:"tools,omitempty"`
	TopP               *float64         `json:"top_p,omitempty"`
	Truncation         string           `json:"truncation,omitempty"`
	Usage              *ResponseUsage   `json:"usage,omitempty"`
	User               *string          `json:"user,omitempty"`
	Metadata           map[string]any   `json:"metadata,omitempty"`
}

// ResponseOutput represents an output item in the response
type ResponseOutput struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Status  string          `json:"status"`
	Role    string          `json:"role,omitempty"`
	Content []OutputContent `json:"content,omitempty"`

	// Function call fields (when Type == "function_call")
	CallID    string `json:"call_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
} // OutputContent represents content within a response output
type OutputContent struct {
	Type        string       `json:"type"`
	Text        string       `json:"text,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

// Annotation represents an annotation in the content
type Annotation struct {
	Type       string `json:"type"`
	StartIndex int    `json:"start_index,omitempty"`
	EndIndex   int    `json:"end_index,omitempty"`
	URL        string `json:"url,omitempty"`
	Title      string `json:"title,omitempty"`
}

// ResponseUsage represents token usage information
type ResponseUsage struct {
	InputTokens         int                    `json:"input_tokens"`
	InputTokensDetails  *ResponseTokensDetails `json:"input_tokens_details,omitempty"`
	OutputTokens        int                    `json:"output_tokens"`
	OutputTokensDetails *ResponseTokensDetails `json:"output_tokens_details,omitempty"`
	TotalTokens         int                    `json:"total_tokens"`
}

// ResponseTokensDetails provides detailed token usage information
type ResponseTokensDetails struct {
	CachedTokens    int `json:"cached_tokens,omitempty"`
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

// ResponseMessage represents a message input for the responses API
type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ResponseToolChoice represents tool choice options
type ResponseToolChoice struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// Tool choice constants
const (
	ResponseToolChoiceAuto     = "auto"
	ResponseToolChoiceRequired = "required"
	ResponseToolChoiceNone     = "none"
)

// ResponseFunctionCallOutput represents function call output for input
type ResponseFunctionCallOutput struct {
	Type   string `json:"type"`    // "function_call_output"
	CallID string `json:"call_id"` // The call_id from function call
	Output string `json:"output"`  // Function execution result
}

// ResponseOptions for creating responses
type ResponseOptions map[string]any

// SetInstructions sets the instructions parameter
func (o ResponseOptions) SetInstructions(instructions string) ResponseOptions {
	o["instructions"] = instructions
	return o
}

// SetMaxOutputTokens sets the max_output_tokens parameter
func (o ResponseOptions) SetMaxOutputTokens(maxTokens int) ResponseOptions {
	o["max_output_tokens"] = maxTokens
	return o
}

// SetTemperature sets the temperature parameter
func (o ResponseOptions) SetTemperature(temperature float64) ResponseOptions {
	o["temperature"] = temperature
	return o
}

// SetTopP sets the top_p parameter
func (o ResponseOptions) SetTopP(topP float64) ResponseOptions {
	o["top_p"] = topP
	return o
}

// SetStore sets the store parameter
func (o ResponseOptions) SetStore(store bool) ResponseOptions {
	o["store"] = store
	return o
}

// SetStream sets the stream parameter with callback
func (o ResponseOptions) SetStream(cb responseCallback) ResponseOptions {
	o["stream"] = cb
	return o
}

// SetUser sets the user parameter
func (o ResponseOptions) SetUser(user string) ResponseOptions {
	o["user"] = user
	return o
}

// SetMetadata sets the metadata parameter
func (o ResponseOptions) SetMetadata(metadata map[string]any) ResponseOptions {
	o["metadata"] = metadata
	return o
}

// SetTools sets the tools parameter
func (o ResponseOptions) SetTools(tools []any) ResponseOptions {
	o["tools"] = tools
	return o
}

// SetToolChoice sets the tool_choice parameter
func (o ResponseOptions) SetToolChoice(toolChoice any) ResponseOptions {
	o["tool_choice"] = toolChoice
	return o
}

// SetToolChoiceAuto sets tool_choice to "auto"
func (o ResponseOptions) SetToolChoiceAuto() ResponseOptions {
	o["tool_choice"] = ResponseToolChoiceAuto
	return o
}

// SetToolChoiceRequired sets tool_choice to "required"
func (o ResponseOptions) SetToolChoiceRequired() ResponseOptions {
	o["tool_choice"] = ResponseToolChoiceRequired
	return o
}

// SetToolChoiceNone sets tool_choice to "none"
func (o ResponseOptions) SetToolChoiceNone() ResponseOptions {
	o["tool_choice"] = ResponseToolChoiceNone
	return o
}

// SetToolChoiceFunction sets tool_choice to force a specific function
func (o ResponseOptions) SetToolChoiceFunction(functionName string) ResponseOptions {
	o["tool_choice"] = ResponseToolChoice{
		Type: "function",
		Name: functionName,
	}
	return o
}

// SetParallelToolCalls sets the parallel_tool_calls parameter
func (o ResponseOptions) SetParallelToolCalls(parallel bool) ResponseOptions {
	o["parallel_tool_calls"] = parallel
	return o
}

// responseCallback defines the callback function for streaming responses
type responseCallback func(response ResponseStreamEvent, done bool, err error)

// ResponseStreamEvent represents a streaming event from the responses API
type ResponseStreamEvent struct {
	Type string `json:"type"`

	// For response events
	Response   *Response `json:"response,omitempty"`
	ResponseID *string   `json:"response_id,omitempty"`

	// For output item events
	OutputIndex *int            `json:"output_index,omitempty"`
	Item        *ResponseOutput `json:"item,omitempty"`

	// For content part events
	ItemID       *string        `json:"item_id,omitempty"`
	ContentIndex *int           `json:"content_index,omitempty"`
	Part         *OutputContent `json:"part,omitempty"`

	// For delta events
	Delta *string `json:"delta,omitempty"`

	// For done events
	Text      *string `json:"text,omitempty"`
	Arguments *string `json:"arguments,omitempty"`
}

// CreateResponse creates a response using the OpenAI Responses API
func (c *Client) CreateResponse(model string, input any, options ResponseOptions) (response Response, err error) {
	if options == nil {
		options = ResponseOptions{}
	}
	options["model"] = model
	options["input"] = input

	if options["stream"] != nil {
		cb := options["stream"].(responseCallback)
		options["stream"] = true
		_, err := c.postCBResponses("v1/responses", options, cb)
		return Response{}, err
	}

	var bytes []byte
	if bytes, err = c.post("v1/responses", options); err == nil {
		if err = json.Unmarshal(bytes, &response); err == nil {
			if response.Error == nil {
				return response, nil
			} else {
				err = response.Error.err()
			}
		}
	}

	return Response{}, err
}

// CreateResponseWithContext creates a response with context support
func (c *Client) CreateResponseWithContext(ctx context.Context, model string, input any, options ResponseOptions) (response Response, err error) {
	if options == nil {
		options = ResponseOptions{}
	}
	options["model"] = model
	options["input"] = input

	if options["stream"] != nil {
		cb := options["stream"].(responseCallback)
		options["stream"] = true
		_, err := c.postCBResponsesWithContext(ctx, "v1/responses", options, cb)
		return Response{}, err
	}

	var bytes []byte
	if bytes, err = c.postWithContext(ctx, "v1/responses", options); err == nil {
		if err = json.Unmarshal(bytes, &response); err == nil {
			if response.Error == nil {
				return response, nil
			} else {
				err = response.Error.err()
			}
		}
	}

	return Response{}, err
}

// CreateResponseStream creates a streaming response
func (c *Client) CreateResponseStream(model string, input any, options ResponseOptions, cb responseCallback) (err error) {
	if options == nil {
		options = ResponseOptions{}
	}
	options["model"] = model
	options["input"] = input
	options["stream"] = true

	_, err = c.postCBResponses("v1/responses", options, cb)
	return err
}

// CreateResponseStreamWithContext creates a streaming response with context support
func (c *Client) CreateResponseStreamWithContext(ctx context.Context, model string, input any, options ResponseOptions, cb responseCallback) (err error) {
	if options == nil {
		options = ResponseOptions{}
	}
	options["model"] = model
	options["input"] = input
	options["stream"] = true

	_, err = c.postCBResponsesWithContext(ctx, "v1/responses", options, cb)
	return err
}

// Helper functions for input creation

// NewResponseMessageInput creates input from a slice of ResponseMessage
func NewResponseMessageInput(messages []ResponseMessage) []ResponseMessage {
	return messages
}

// NewResponseTextInput creates input from a text string
func NewResponseTextInput(text string) string {
	return text
}

// NewResponseMessage creates a new ResponseMessage
func NewResponseMessage(role, content string) ResponseMessage {
	return ResponseMessage{
		Role:    role,
		Content: content,
	}
}

// NewResponseFunctionCallOutput creates a function call output for input
func NewResponseFunctionCallOutput(callID, output string) ResponseFunctionCallOutput {
	return ResponseFunctionCallOutput{
		Type:   "function_call_output",
		CallID: callID,
		Output: output,
	}
}

// NewResponseTool creates a function tool for responses API (reuses existing Tool from assistants)
func NewResponseTool(name, description string, parameters ToolFunctionParameters) Tool {
	return NewFunctionTool(ToolFunction{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	})
}

// ArgumentsParsed returns the parsed arguments from a function call ResponseOutput
func (r ResponseOutput) ArgumentsParsed() (result map[string]any, err error) {
	if r.Type == "function_call" && r.Arguments != "" {
		err = json.Unmarshal([]byte(r.Arguments), &result)
	}
	return result, err
}

// ArgumentsInto parses the function call arguments into a given interface
func (r ResponseOutput) ArgumentsInto(out any) (err error) {
	if r.Type != "function_call" {
		return fmt.Errorf("not a function call output")
	}
	if r.Arguments == "" {
		return fmt.Errorf("arguments are empty")
	}
	return json.Unmarshal([]byte(r.Arguments), &out)
}
