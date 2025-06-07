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
