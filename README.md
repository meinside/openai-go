# openai-go

OpenAI API wrapper library for Go.

## How to use

Generate API key from [here](https://platform.openai.com/account/api-keys),

and get organization id from [here](https://platform.openai.com/account/org-settings).

```go
const apiKey = "ab-cdefHIJKLMNOPQRSTUVWXYZ0123456789"
const orgID = "org-0123456789ABCDefghijklmnopQRSTUVWxyz"

func main() {
    client := NewClient(apiKey, orgID)

    if models, err := client.ListModels(); err != nil {
        log.Printf("available models = %+v", models.Data)
    }
}
```

## How to test

Export following environment variables:

```bash
$ export OPENAI_API_KEY=ab-cdefHIJKLMNOPQRSTUVWXYZ0123456789
$ export OPENAI_ORGANIZATION=org-0123456789ABCDefghijklmnopQRSTUVWxyz

# for verbose messages and http request dumps,
$ export VERBOSE=true

```

and then

```bash
$ go test
```

**CAUTION**: It is advised to [set usage limits](https://platform.openai.com/account/limits) before running tests; running all tests at once costs about ~$0.2.

## Todos/WIP

### Implemented

All API functions so far (2023.11.07.) are implemented, but not all of them were tested on a paid account.

- [X] [Audio](https://platform.openai.com/docs/api-reference/audio)
- [X] [Chat](https://platform.openai.com/docs/api-reference/chat)
- [X] [Completions](https://platform.openai.com/docs/api-reference/completions)
- [X] [Embeddings](https://platform.openai.com/docs/api-reference/embeddings)
- [X] [Fine-tuning](https://platform.openai.com/docs/api-reference/fine-tuning)
- [X] [Files](https://platform.openai.com/docs/api-reference/files)
- [X] [Images](https://platform.openai.com/docs/api-reference/images)
- [X] [Models](https://platform.openai.com/docs/api-reference/models): works on a non-paid account
- [X] [Moderations](https://platform.openai.com/docs/api-reference/moderations): works on a non-paid account
- [X] [Responses](https://platform.openai.com/docs/api-reference/responses)

### Responses API examples

```go
// Simple text input
response, err := client.CreateResponse("gpt-4.1", "Hello, how can you help me?", nil)
if err != nil {
    log.Fatal(err)
}
log.Printf("Response: %s", response.Output[0].Content[0].Text)

// With message array input
messages := []openai.ResponseMessage{
    openai.NewResponseMessage("user", "What is the weather like?"),
    openai.NewResponseMessage("assistant", "I'd be happy to help with weather information. Could you please specify your location?"),
    openai.NewResponseMessage("user", "New York City"),
}

response, err := client.CreateResponse("gpt-4.1", messages, nil)
if err != nil {
    log.Fatal(err)
}

// With options
options := openai.ResponseOptions{}
options.SetInstructions("You are a helpful assistant.")
options.SetTemperature(0.7)
options.SetMaxOutputTokens(100)

response, err := client.CreateResponse("gpt-4.1", "Explain quantum computing", options)
if err != nil {
    log.Fatal(err)
}
```

#### Using Tools (Function Calling)

```go
// Create a weather tool
weatherTool := openai.NewResponseTool("get_weather", 
    "Get current temperature for a given location.", 
    openai.NewToolFunctionParameters().
        AddPropertyWithDescription("location", "string", "City and country e.g. Bogotá, Colombia").
        SetRequiredParameters([]string{"location"}))

// Configure options with tools
options := openai.ResponseOptions{}
options.SetInstructions("You are a helpful weather assistant.")
options.SetTools([]any{weatherTool})
options.SetToolChoiceAuto() // or SetToolChoiceRequired(), SetToolChoiceFunction("get_weather")

// First call - model decides to call function
response, err := client.CreateResponse("gpt-4o", "What's the weather in Paris?", options)
if err != nil {
    log.Fatal(err)
}

// Check for function calls
for _, output := range response.Output {
    if output.Type == "function_call" {
        log.Printf("Function call: %s", output.Name)
        
        // Parse arguments
        args, err := output.ArgumentsParsed()
        if err != nil {
            log.Fatal(err)
        }
        
        // Execute your function (simulate)
        result := "22°C, sunny"
        
        // Create input with function result
        input := []any{
            openai.NewResponseMessage("user", "What's the weather in Paris?"),
            output, // Include the function call
            openai.NewResponseFunctionCallOutput(output.CallID, result),
        }
        
        // Second call - get final response with function result
        finalResponse, err := client.CreateResponse("gpt-4o", input, openai.ResponseOptions{})
        if err != nil {
            log.Fatal(err)
        }
        
        log.Printf("Final response: %s", finalResponse.Output[0].Content[0].Text)
        break
    }
}
```

#### Streaming Responses

```go
err := client.CreateResponseStream("gpt-4.1", "Tell me a story", nil, func(event openai.ResponseStreamEvent, done bool, err error) {
    if err != nil {
        log.Printf("Stream error: %v", err)
        return
    }
    
    switch event.Type {
    case "response.output_text.delta":
        if event.Delta != nil {
            fmt.Print(*event.Delta)
        }
    case "response.completed":
        fmt.Println("\nStream completed")
    }
    
    if done {
        return
    }
})

if err != nil {
    log.Fatal(err)
}
```

#### Streaming with Tools

```go
weatherTool := openai.NewResponseTool("get_weather", "Get weather info", 
    openai.NewToolFunctionParameters().
        AddPropertyWithDescription("location", "string", "City name").
        SetRequiredParameters([]string{"location"}))

options := openai.ResponseOptions{}
options.SetTools([]any{weatherTool})
options.SetToolChoiceAuto()

err := client.CreateResponseStream("gpt-4o", "Weather in Tokyo?", options, 
    func(event openai.ResponseStreamEvent, done bool, err error) {
        if err != nil {
            log.Printf("Stream error: %v", err)
            return
        }
        
        switch event.Type {
        case "response.output_item.added":
            if event.Item != nil && event.Item.Type == "function_call" {
                log.Printf("Function call started: %s", event.Item.Name)
            }
        case "response.function_call_arguments.delta":
            if event.Delta != nil {
                fmt.Print(*event.Delta) // Print argument deltas
            }
        case "response.function_call_arguments.done":
            if event.Arguments != nil {
                log.Printf("\nFunction arguments complete: %s", *event.Arguments)
            }
        case "response.completed":
            log.Println("Stream completed")
        }
        
        if done {
            return
        }
    })

if err != nil {
    log.Fatal(err)
}
```

### Beta

- [X] [Assistants](https://platform.openai.com/docs/api-reference/assistants)
- [X] [Threads](https://platform.openai.com/docs/api-reference/threads)
- [X] [Messages](https://platform.openai.com/docs/api-reference/messages)
- [X] [Runs](https://platform.openai.com/docs/api-reference/runs)

#### Note

Beta API functions require beta header like this:

```go
client.SetBetaHeader(`assistants=v1`)
```

### Help Wanted

- [X] ~~Stream([server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)) options are not implemented yet.~~ thanks to @tectiv3 :-)
- [ ] Add some sample applications.

