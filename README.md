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

    if models, err := client.LitModels(); err != nil {
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

**CAUTION**: It is advised to set usage limits before running tests; running all tests at once costs about ~$0.2.

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

- [X] ~~Stream([server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)) options are not implemented yet.~~ thanks to @tectiv3 :-)

### Beta

- [ ] [Assistants](https://platform.openai.com/docs/api-reference/assistants)
- [ ] [Threads](https://platform.openai.com/docs/api-reference/threads)
- [ ] [Messages](https://platform.openai.com/docs/api-reference/messages)
- [ ] [Runs](https://platform.openai.com/docs/api-reference/runs)

## License

MIT

