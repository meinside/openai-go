package openai

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	responsesModel = "gpt-4o"
)

// TestResponsesReal tests the Responses API with real API calls
func TestResponsesReal(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
		return
	}

	// Test basic text input
	response, err := client.CreateResponse(responsesModel, "Hello! How can you help me?", nil)
	if err != nil {
		t.Errorf("CreateResponse failed: %v", err)
		return
	}

	if response.ID == "" {
		t.Error("Expected response to have an ID")
	}

	if response.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", response.Status)
	}

	if len(response.Output) == 0 {
		t.Error("Expected response to have output")
	}

	log.Printf("Real API response ID: %s", response.ID)
	if len(response.Output) > 0 && len(response.Output[0].Content) > 0 {
		log.Printf("Real API response text: %s", response.Output[0].Content[0].Text)
	}
}

// TestResponsesRealWithMessages tests with message array input
func TestResponsesRealWithMessages(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
		return
	}

	// Test with message array input
	messages := []ResponseMessage{
		NewResponseMessage("user", "knock knock."),
		NewResponseMessage("assistant", "Who's there?"),
		NewResponseMessage("user", "Orange."),
	}

	response, err := client.CreateResponse(responsesModel, messages, nil)
	if err != nil {
		t.Errorf("CreateResponse with messages failed: %v", err)
		return
	}

	if response.ID == "" {
		t.Error("Expected response to have an ID")
	}

	log.Printf("Real API response with messages ID: %s", response.ID)
	if len(response.Output) > 0 && len(response.Output[0].Content) > 0 {
		log.Printf("Real API response text: %s", response.Output[0].Content[0].Text)
	}
}

// TestResponsesRealWithOptions tests with various options
func TestResponsesRealWithOptions(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
		return
	}

	// Test with options
	options := ResponseOptions{}
	options.SetInstructions("You are a helpful assistant. Keep responses brief.")
	options.SetTemperature(0.7)
	options.SetMaxOutputTokens(50)

	response, err := client.CreateResponse(responsesModel, "Explain quantum computing in one sentence.", options)
	if err != nil {
		t.Errorf("CreateResponse with options failed: %v", err)
		return
	}

	if response.ID == "" {
		t.Error("Expected response to have an ID")
	}

	log.Printf("Real API response with options ID: %s", response.ID)
	if len(response.Output) > 0 && len(response.Output[0].Content) > 0 {
		log.Printf("Real API response text: %s", response.Output[0].Content[0].Text)
	}
}

// TestResponsesRealStream tests streaming functionality
func TestResponsesRealStream(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
		return
	}

	events := []ResponseStreamEvent{}
	callbackCount := 0
	done := make(chan bool)

	// Test streaming
	err := client.CreateResponseStream(responsesModel, "Tell me a short joke", nil, func(event ResponseStreamEvent, isDone bool, err error) {
		callbackCount++
		if err != nil {
			t.Errorf("Stream callback error: %v", err)
			done <- true
			return
		}

		events = append(events, event)
		log.Printf("Stream event type: %s", event.Type)

		// Print delta text if available
		if event.Delta != nil {
			log.Printf("Delta: %s", *event.Delta)
		}

		if isDone {
			log.Printf("Stream completed with %d events", len(events))
			done <- true
		}
	})
	if err != nil {
		t.Errorf("CreateResponseStream failed: %v", err)
		return
	}

	// Wait for streaming to complete (with timeout)
	select {
	case <-done:
		// Success
	case <-time.After(30 * time.Second):
		t.Error("Stream test timed out")
	}

	if callbackCount == 0 {
		t.Error("Expected stream callbacks to be called, but they weren't")
	}

	if len(events) == 0 {
		t.Error("Expected streaming events, got none")
	}
}

// Mock tests below - these run without needing API keys

func TestCreateResponseMock(t *testing.T) {
	// Mock response JSON
	mockResponse := `{
		"id": "resp_67ccd3a9da748190baa7f1570fe91ac604becb25c45c1d41",
		"object": "response",
		"created_at": 1741476777,
		"status": "completed",
		"error": null,
		"model": "gpt-4o-2024-08-06",
		"output": [
			{
				"type": "message",
				"id": "msg_67ccd3acc8d48190a77525dc6de64b4104becb25c45c1d41",
				"status": "completed",
				"role": "assistant",
				"content": [
					{
						"type": "output_text",
						"text": "Hello! How can I help you today?",
						"annotations": []
					}
				]
			}
		],
		"usage": {
			"input_tokens": 328,
			"output_tokens": 52,
			"total_tokens": 380
		}
	}`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/responses" {
			t.Errorf("Expected path /v1/responses, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check request body
		var requestBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if requestBody["model"] != "gpt-4.1" {
			t.Errorf("Expected model gpt-4.1, got %v", requestBody["model"])
		}

		if requestBody["input"] != "Hello!" {
			t.Errorf("Expected input 'Hello!', got %v", requestBody["input"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-key", "test-org")
	client.baseURL = &server.URL

	// Test CreateResponse
	response, err := client.CreateResponse("gpt-4.1", "Hello!", nil)
	if err != nil {
		t.Errorf("CreateResponse failed: %v", err)
	}

	if response.ID != "resp_67ccd3a9da748190baa7f1570fe91ac604becb25c45c1d41" {
		t.Errorf("Expected response ID resp_67ccd3a9da748190baa7f1570fe91ac604becb25c45c1d41, got %s", response.ID)
	}

	if response.Status != "completed" {
		t.Errorf("Expected status completed, got %s", response.Status)
	}

	if len(response.Output) != 1 {
		t.Errorf("Expected 1 output item, got %d", len(response.Output))
	}

	if len(response.Output) > 0 && len(response.Output[0].Content) > 0 {
		if response.Output[0].Content[0].Text != "Hello! How can I help you today?" {
			t.Errorf("Expected text 'Hello! How can I help you today?', got %s", response.Output[0].Content[0].Text)
		}
	}
}

func TestCreateResponseWithMessagesMock(t *testing.T) {
	// Mock response JSON
	mockResponse := `{
		"id": "resp_test123",
		"object": "response",
		"created_at": 1741476777,
		"status": "completed",
		"model": "gpt-4.1",
		"output": [
			{
				"type": "message",
				"id": "msg_test123",
				"status": "completed",
				"role": "assistant",
				"content": [
					{
						"type": "output_text",
						"text": "Orange you glad I didn't say banana!",
						"annotations": []
					}
				]
			}
		]
	}`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request body
		var requestBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify input is an array
		input, ok := requestBody["input"].([]any)
		if !ok {
			t.Errorf("Expected input to be an array, got %T", requestBody["input"])
		}

		if len(input) != 3 {
			t.Errorf("Expected 3 input messages, got %d", len(input))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-key", "test-org")
	client.baseURL = &server.URL

	// Test with message array input
	messages := []ResponseMessage{
		NewResponseMessage("user", "knock knock."),
		NewResponseMessage("assistant", "Who's there?"),
		NewResponseMessage("user", "Orange."),
	}

	response, err := client.CreateResponse("gpt-4.1", messages, nil)
	if err != nil {
		t.Errorf("CreateResponse failed: %v", err)
	}

	if response.ID != "resp_test123" {
		t.Errorf("Expected response ID resp_test123, got %s", response.ID)
	}
}

func TestCreateResponseWithOptionsMock(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request body
		var requestBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify options are present
		if requestBody["instructions"] != "You are a helpful assistant." {
			t.Errorf("Expected instructions 'You are a helpful assistant.', got %v", requestBody["instructions"])
		}

		if requestBody["temperature"] != 0.7 {
			t.Errorf("Expected temperature 0.7, got %v", requestBody["temperature"])
		}

		if requestBody["max_output_tokens"] != float64(100) {
			t.Errorf("Expected max_output_tokens 100, got %v", requestBody["max_output_tokens"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"test","object":"response","created_at":1234567890,"status":"completed","model":"gpt-4.1","output":[]}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-key", "test-org")
	client.baseURL = &server.URL

	// Test with options
	options := ResponseOptions{}
	options.SetInstructions("You are a helpful assistant.")
	options.SetTemperature(0.7)
	options.SetMaxOutputTokens(100)

	_, err := client.CreateResponse("gpt-4.1", "Hello!", options)
	if err != nil {
		t.Errorf("CreateResponse failed: %v", err)
	}
}

func TestCreateResponseWithContextMock(t *testing.T) {
	// Create test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"test","object":"response","created_at":1234567890,"status":"completed","model":"gpt-4.1","output":[]}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-key", "test-org")
	client.baseURL = &server.URL

	// Test context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.CreateResponseWithContext(ctx, "gpt-4.1", "Hello!", nil)
	if err == nil {
		t.Error("Expected context timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected context deadline exceeded error, got: %v", err)
	}
}

func TestCreateResponseStreamMock(t *testing.T) {
	// Mock streaming response
	streamingResponse := `event: response.created
data: {"type":"response.created","response":{"id":"resp_test","object":"response","created_at":1741290958,"status":"in_progress","model":"gpt-4.1","output":[]}}

event: response.output_text.delta
data: {"type":"response.output_text.delta","item_id":"msg_test","output_index":0,"content_index":0,"delta":"Hi"}

event: response.completed
data: {"type":"response.completed","response":{"id":"resp_test","object":"response","created_at":1741290958,"status":"completed","model":"gpt-4.1","output":[{"type":"message","id":"msg_test","status":"completed","role":"assistant","content":[{"type":"output_text","text":"Hi there!","annotations":[]}]}]}}

`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for streaming request
		var requestBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if requestBody["stream"] != true {
			t.Errorf("Expected stream to be true, got %v", requestBody["stream"])
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(streamingResponse))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-key", "test-org")
	client.baseURL = &server.URL

	events := []ResponseStreamEvent{}
	callbackCount := 0

	// Test streaming
	err := client.CreateResponseStream("gpt-4.1", "Hello!", nil, func(event ResponseStreamEvent, done bool, err error) {
		callbackCount++
		if err != nil {
			t.Errorf("Callback error: %v", err)
			return
		}
		events = append(events, event)
		if done {
			return
		}
	})
	if err != nil {
		t.Errorf("CreateResponseStream failed: %v", err)
	}

	// Give some time for streaming to complete
	time.Sleep(100 * time.Millisecond)

	if callbackCount == 0 {
		t.Error("Expected callback to be called, but it wasn't")
	}

	if len(events) == 0 {
		t.Error("Expected streaming events, got none")
	}
}

func TestResponseOptionSetters(t *testing.T) {
	options := ResponseOptions{}

	// Test all setter methods
	options.SetInstructions("test instructions")
	options.SetMaxOutputTokens(500)
	options.SetTemperature(0.8)
	options.SetTopP(0.9)
	options.SetStore(true)
	options.SetUser("user123")
	options.SetMetadata(map[string]any{"key": "value"})
	options.SetTools([]any{})
	options.SetToolChoice("auto")
	options.SetParallelToolCalls(true)

	// Verify all options are set
	if options["instructions"] != "test instructions" {
		t.Errorf("Expected instructions 'test instructions', got %v", options["instructions"])
	}

	if options["max_output_tokens"] != 500 {
		t.Errorf("Expected max_output_tokens 500, got %v", options["max_output_tokens"])
	}

	if options["temperature"] != 0.8 {
		t.Errorf("Expected temperature 0.8, got %v", options["temperature"])
	}

	if options["top_p"] != 0.9 {
		t.Errorf("Expected top_p 0.9, got %v", options["top_p"])
	}

	if options["store"] != true {
		t.Errorf("Expected store true, got %v", options["store"])
	}

	if options["user"] != "user123" {
		t.Errorf("Expected user 'user123', got %v", options["user"])
	}

	metadata, ok := options["metadata"].(map[string]any)
	if !ok || metadata["key"] != "value" {
		t.Errorf("Expected metadata with key='value', got %v", options["metadata"])
	}

	if options["tool_choice"] != "auto" {
		t.Errorf("Expected tool_choice 'auto', got %v", options["tool_choice"])
	}

	if options["parallel_tool_calls"] != true {
		t.Errorf("Expected parallel_tool_calls true, got %v", options["parallel_tool_calls"])
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test NewResponseMessage
	msg := NewResponseMessage("user", "Hello")
	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got %s", msg.Role)
	}
	if msg.Content != "Hello" {
		t.Errorf("Expected content 'Hello', got %s", msg.Content)
	}

	// Test NewResponseTextInput
	textInput := NewResponseTextInput("Hello world")
	if textInput != "Hello world" {
		t.Errorf("Expected 'Hello world', got %s", textInput)
	}

	// Test NewResponseMessageInput
	messages := []ResponseMessage{
		NewResponseMessage("user", "Hello"),
		NewResponseMessage("assistant", "Hi there!"),
	}
	messageInput := NewResponseMessageInput(messages)
	if len(messageInput) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messageInput))
	}
	if messageInput[0].Role != "user" {
		t.Errorf("Expected first message role 'user', got %s", messageInput[0].Role)
	}
	if messageInput[1].Role != "assistant" {
		t.Errorf("Expected second message role 'assistant', got %s", messageInput[1].Role)
	}
}

// Tests for Tools functionality

func TestResponsesRealWithTools(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
		return
	}

	// Create a search tool
	searchTool := NewBuiltinTool("web_search_preview")

	// Test with tools
	options := ResponseOptions{}
	options.SetInstructions("You are a helpful assistant.")
	options.SetTools([]any{searchTool})
	options.SetToolChoiceAuto()

	response, err := client.CreateResponse(responsesModel, "What was a positive news story from today?", options)
	if err != nil {
		t.Errorf("CreateResponse with tools failed: %v", err)
		return
	}

	if response.ID == "" {
		t.Error("Expected response to have an ID")
	}

	log.Printf("Real API response with tools ID: %s", response.ID)

	// Check if function call was made
	functionCalled := false
	for _, output := range response.Output {
		if output.Type == "web_search_call" {
			functionCalled = true
			log.Printf("Search call with status: %s", output.Status)
		} else if output.Type == "message" {
			if len(output.Content) > 0 && len(output.Content) > 0 {
				log.Printf("Final response text: %s", output.Content[0].Text)
				log.Printf("Annotations: %v", output.Content[0].Annotations)
			}
		} else {
			t.Errorf("Unexpected output type: %s", output.Type)
		}
	}

	if !functionCalled {
		log.Printf("Note: Function was not called by the model for this request")
	}
}

func TestResponsesRealToolsWithCallback(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
		return
	}

	// Create a simple calculation tool
	calcTool := NewResponseTool("calculate", "Perform basic math calculations.",
		NewToolFunctionParameters().
			AddPropertyWithDescription("expression", "string", "Math expression to evaluate").
			SetRequiredParameters([]string{"expression"}))

	// First call - get function call
	options := ResponseOptions{}
	options.SetTools([]any{calcTool})
	options.SetToolChoiceRequired()

	response, err := client.CreateResponse(responsesModel, "Calculate 25 * 4", options)
	if err != nil {
		t.Errorf("CreateResponse with required tools failed: %v", err)
		return
	}

	// Find function call
	var functionCall *ResponseOutput
	for i, output := range response.Output {
		if output.Type == "function_call" {
			functionCall = &response.Output[i]
			break
		}
	}

	if functionCall == nil {
		t.Error("Expected function call in response")
		return
	}

	log.Printf("Function call received: %s(%s)", functionCall.Name, functionCall.Arguments)

	// Simulate function execution
	functionResult := "100"

	// Create input with original messages + function call + result
	input := []any{
		NewResponseMessage("user", "Calculate 25 * 4"),
		*functionCall, // Add the function call
		NewResponseFunctionCallOutput(functionCall.CallID, functionResult), // Add the result
	}

	// Second call - get final response with function result
	finalResponse, err := client.CreateResponse(responsesModel, input, ResponseOptions{})
	if err != nil {
		t.Errorf("CreateResponse with function result failed: %v", err)
		return
	}

	log.Printf("Final response ID: %s", finalResponse.ID)
	if len(finalResponse.Output) > 0 && len(finalResponse.Output[0].Content) > 0 {
		log.Printf("Final response text: %s", finalResponse.Output[0].Content[0].Text)
	}
}

func TestResponsesRealStreamWithTools(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
		return
	}

	// Create a simple tool
	timeTool := NewResponseTool("get_current_time", "Get the current time.",
		NewToolFunctionParameters().
			AddPropertyWithDescription("timezone", "string", "Timezone (optional)"))
	log.Printf("Using tool: %v", timeTool)
	options := ResponseOptions{}
	options.SetTools([]any{timeTool})
	options.SetToolChoiceAuto()

	events := []ResponseStreamEvent{}
	callbackCount := 0
	done := make(chan bool)
	functionCallEvents := 0
	var completedFunctionCall *ResponseOutput

	// Test streaming with tools
	err := client.CreateResponseStream(responsesModel, "What time is it?", options, func(event ResponseStreamEvent, isDone bool, err error) {
		callbackCount++
		if err != nil {
			t.Errorf("Stream callback error: %v", err)
			done <- true
			return
		}

		events = append(events, event)
		log.Printf("Stream event type: %s", event.Type)

		// Handle different event types
		switch event.Type {
		case "response.output_item.added":
			if event.Item != nil && event.Item.Type == "function_call" {
				log.Printf("Function call started: %s", event.Item.Name)
			}
		case "response.function_call_arguments.delta":
			functionCallEvents++
			if event.Delta != nil {
				log.Printf("Function call arguments delta: %s", *event.Delta)
			}
		case "response.function_call_arguments.done":
			functionCallEvents++
			if event.Arguments != nil {
				log.Printf("Function call arguments completed: %s", *event.Arguments)
			}
		case "response.output_item.done":
			if event.Item != nil && event.Item.Type == "function_call" {
				completedFunctionCall = event.Item
				log.Printf("\n=== FUNCTION CALL COMPLETED ===")
				log.Printf("Function: %s", event.Item.Name)
				log.Printf("Call ID: %s", event.Item.CallID)
				log.Printf("Raw Arguments: %s", event.Item.Arguments)

				// Parse arguments using ArgumentsParsed
				if args, parseErr := event.Item.ArgumentsParsed(); parseErr == nil {
					log.Printf("Parsed Arguments: %+v", args)
					for key, value := range args {
						log.Printf("  %s: %v (type: %T)", key, value, value)
					}
				} else {
					log.Printf("Failed to parse arguments: %v", parseErr)
				}
				log.Printf("===============================\n")
			}
		case "response.completed":
			log.Printf("Response stream completed")
		}

		if isDone {
			log.Printf("Stream completed with %d events, %d function call events", len(events), functionCallEvents)
			if completedFunctionCall != nil {
				log.Printf("Final function call summary: %s with args %s", completedFunctionCall.Name, completedFunctionCall.Arguments)
			}
			done <- true
		}
	})
	if err != nil {
		t.Errorf("CreateResponseStream with tools failed: %v", err)
		return
	}

	// Wait for streaming to complete (with timeout)
	select {
	case <-done:
		// Success
	case <-time.After(30 * time.Second):
		t.Error("Stream test with tools timed out")
	}

	if callbackCount == 0 {
		t.Error("Expected stream callbacks to be called, but they weren't")
	}
}

// Mock tests for tools functionality

func TestCreateResponseWithToolsMock(t *testing.T) {
	// Mock response with function call
	mockResponse := `{
		"id": "resp_tools_test",
		"object": "response",
		"created_at": 1741476777,
		"status": "completed",
		"model": "gpt-4o",
		"output": [
			{
				"type": "function_call",
				"id": "fc_12345xyz",
				"call_id": "call_12345xyz",
				"name": "get_weather",
				"arguments": "{\"location\":\"Paris, France\"}"
			}
		]
	}`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request body
		var requestBody map[string]any
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify tools are present
		tools, ok := requestBody["tools"].([]any)
		if !ok {
			t.Errorf("Expected tools to be an array, got %T", requestBody["tools"])
		}

		if len(tools) != 1 {
			t.Errorf("Expected 1 tool, got %d", len(tools))
		}

		// Verify tool_choice
		if requestBody["tool_choice"] != "auto" {
			t.Errorf("Expected tool_choice 'auto', got %v", requestBody["tool_choice"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-key", "test-org")
	client.baseURL = &server.URL

	// Create tool
	weatherTool := NewResponseTool("get_weather", "Get weather for location",
		NewToolFunctionParameters().
			AddPropertyWithDescription("location", "string", "City and country").
			SetRequiredParameters([]string{"location"}))

	// Test with tools
	options := ResponseOptions{}
	options.SetTools([]any{weatherTool})
	options.SetToolChoiceAuto()

	response, err := client.CreateResponse("gpt-4o", "What is the weather in Paris?", options)
	if err != nil {
		t.Errorf("CreateResponse with tools failed: %v", err)
	}

	if len(response.Output) != 1 {
		t.Errorf("Expected 1 output item, got %d", len(response.Output))
	}

	if response.Output[0].Type != "function_call" {
		t.Errorf("Expected function_call type, got %s", response.Output[0].Type)
	}

	if response.Output[0].Name != "get_weather" {
		t.Errorf("Expected function name 'get_weather', got %s", response.Output[0].Name)
	}

	// Test argument parsing
	args, err := response.Output[0].ArgumentsParsed()
	if err != nil {
		t.Errorf("Failed to parse arguments: %v", err)
	}

	if args["location"] != "Paris, France" {
		t.Errorf("Expected location 'Paris, France', got %v", args["location"])
	}
}

func TestToolChoiceOptions(t *testing.T) {
	options := ResponseOptions{}

	// Test tool choice setters
	options.SetToolChoiceAuto()
	if options["tool_choice"] != "auto" {
		t.Errorf("Expected tool_choice 'auto', got %v", options["tool_choice"])
	}

	options.SetToolChoiceRequired()
	if options["tool_choice"] != "required" {
		t.Errorf("Expected tool_choice 'required', got %v", options["tool_choice"])
	}

	options.SetToolChoiceNone()
	if options["tool_choice"] != "none" {
		t.Errorf("Expected tool_choice 'none', got %v", options["tool_choice"])
	}

	options.SetToolChoiceFunction("get_weather")
	toolChoice, ok := options["tool_choice"].(ResponseToolChoice)
	if !ok {
		t.Errorf("Expected ResponseToolChoice type, got %T", options["tool_choice"])
	}
	if toolChoice.Type != "function" {
		t.Errorf("Expected type 'function', got %s", toolChoice.Type)
	}
	if toolChoice.Name != "get_weather" {
		t.Errorf("Expected name 'get_weather', got %s", toolChoice.Name)
	}
}

func TestFunctionCallOutput(t *testing.T) {
	// Test creating function call output
	output := NewResponseFunctionCallOutput("call_123", "The weather in Paris is 22°C")
	if output.Type != "function_call_output" {
		t.Errorf("Expected type 'function_call_output', got %s", output.Type)
	}
	if output.CallID != "call_123" {
		t.Errorf("Expected call_id 'call_123', got %s", output.CallID)
	}
	if output.Output != "The weather in Paris is 22°C" {
		t.Errorf("Expected output 'The weather in Paris is 22°C', got %s", output.Output)
	}
}

func TestResponseOutputArgumentParsing(t *testing.T) {
	// Test parsing function call arguments
	output := ResponseOutput{
		Type:      "function_call",
		Name:      "get_weather",
		Arguments: `{"location":"New York","units":"celsius"}`,
	}

	// Test ArgumentsParsed
	args, err := output.ArgumentsParsed()
	if err != nil {
		t.Errorf("ArgumentsParsed failed: %v", err)
	}

	if args["location"] != "New York" {
		t.Errorf("Expected location 'New York', got %v", args["location"])
	}
	if args["units"] != "celsius" {
		t.Errorf("Expected units 'celsius', got %v", args["units"])
	}

	// Test ArgumentsInto with struct
	type WeatherArgs struct {
		Location string `json:"location"`
		Units    string `json:"units"`
	}

	var weatherArgs WeatherArgs
	err = output.ArgumentsInto(&weatherArgs)
	if err != nil {
		t.Errorf("ArgumentsInto failed: %v", err)
	}

	if weatherArgs.Location != "New York" {
		t.Errorf("Expected location 'New York', got %s", weatherArgs.Location)
	}
	if weatherArgs.Units != "celsius" {
		t.Errorf("Expected units 'celsius', got %s", weatherArgs.Units)
	}

	// Test with non-function-call type
	nonFunctionOutput := ResponseOutput{Type: "message"}
	err = nonFunctionOutput.ArgumentsInto(&weatherArgs)
	if err == nil {
		t.Error("Expected error for non-function-call type")
	}

	// Test with empty arguments
	emptyArgsOutput := ResponseOutput{Type: "function_call", Arguments: ""}
	err = emptyArgsOutput.ArgumentsInto(&weatherArgs)
	if err == nil {
		t.Error("Expected error for empty arguments")
	}
}

func TestStreamingFunctionCallEvents(t *testing.T) {
	// Mock streaming response with function calls
	streamingResponse := `event: response.created
data: {"type":"response.created","response":{"id":"resp_test","object":"response","created_at":1741290958,"status":"in_progress","model":"gpt-4o","output":[]}}

event: response.output_item.added
data: {"type":"response.output_item.added","response_id":"resp_test","output_index":0,"item":{"type":"function_call","id":"fc_123","call_id":"call_123","name":"get_weather","arguments":""}}

event: response.function_call_arguments.delta
data: {"type":"response.function_call_arguments.delta","response_id":"resp_test","item_id":"fc_123","output_index":0,"delta":"{\""}

event: response.function_call_arguments.delta
data: {"type":"response.function_call_arguments.delta","response_id":"resp_test","item_id":"fc_123","output_index":0,"delta":"location"}

event: response.function_call_arguments.done
data: {"type":"response.function_call_arguments.done","response_id":"resp_test","item_id":"fc_123","output_index":0,"arguments":"{\"location\":\"Paris\"}"}

event: response.completed
data: {"type":"response.completed","response":{"id":"resp_test","object":"response","created_at":1741290958,"status":"completed","model":"gpt-4o","output":[{"type":"function_call","id":"fc_123","call_id":"call_123","name":"get_weather","arguments":"{\"location\":\"Paris\"}"}]}}

`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(streamingResponse))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-key", "test-org")
	client.baseURL = &server.URL

	events := []ResponseStreamEvent{}
	functionCallEvents := 0

	// Test streaming
	err := client.CreateResponseStream("gpt-4o", "Get weather", nil, func(event ResponseStreamEvent, done bool, err error) {
		if err != nil {
			t.Errorf("Callback error: %v", err)
			return
		}
		events = append(events, event)

		// Count function call related events
		if strings.Contains(event.Type, "function_call") {
			functionCallEvents++
		}

		if done {
			return
		}
	})
	if err != nil {
		t.Errorf("CreateResponseStream failed: %v", err)
	}

	// Give some time for streaming to complete
	time.Sleep(100 * time.Millisecond)

	if functionCallEvents == 0 {
		t.Error("Expected function call events, got none")
	}

	// Verify we got the expected function call events
	expectedEvents := []string{"response.function_call_arguments.delta", "response.function_call_arguments.done"}
	foundEvents := make(map[string]bool)

	for _, event := range events {
		for _, expected := range expectedEvents {
			if event.Type == expected {
				foundEvents[expected] = true
			}
		}
	}

	for _, expected := range expectedEvents {
		if !foundEvents[expected] {
			t.Errorf("Expected to find event type %s, but didn't", expected)
		}
	}
}
