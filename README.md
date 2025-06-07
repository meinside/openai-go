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

