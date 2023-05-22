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

## Todos/WIP

### Implemented

All API functions so far (2023.03.06.) are implemented, but not all of them were tested on a paid account.

- [X] [Models](https://platform.openai.com/docs/api-reference/models): works on a non-paid account
- [X] [Completions](https://platform.openai.com/docs/api-reference/completions)
- [X] [Chat](https://platform.openai.com/docs/api-reference/chat)
- [X] [Edits](https://platform.openai.com/docs/api-reference/edits)
- [X] [Images](https://platform.openai.com/docs/api-reference/images)
- [X] [Embeddings](https://platform.openai.com/docs/api-reference/embeddings)
- [X] [Audio](https://platform.openai.com/docs/api-reference/audio)
- [X] [Files](https://platform.openai.com/docs/api-reference/files)
- [X] [Fine-tunes](https://platform.openai.com/docs/api-reference/fine-tunes)
- [X] [Moderations](https://platform.openai.com/docs/api-reference/moderations): works on a non-paid account

- [X] ~~Stream([server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)) options are not implemented yet.~~ thanks to @tectiv3 :-)

## License

MIT

