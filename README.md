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

All API functions so far are implemented, but not yet tested on a paid account.

- [X] [Models](https://platform.openai.com/docs/api-reference/models): works on a non-paid account
- [ ] [Completions](https://platform.openai.com/docs/api-reference/completions)
- [ ] [Chat](https://platform.openai.com/docs/api-reference/chat)
- [ ] [Edits](https://platform.openai.com/docs/api-reference/edits)
- [ ] [Images](https://platform.openai.com/docs/api-reference/images)
- [ ] [Embeddings](https://platform.openai.com/docs/api-reference/embeddings)
- [ ] [Audio](https://platform.openai.com/docs/api-reference/audio)
- [ ] [Files](https://platform.openai.com/docs/api-reference/files)
- [ ] [Fine-tunes](https://platform.openai.com/docs/api-reference/fine-tunes)
- [X] [Moderations](https://platform.openai.com/docs/api-reference/moderations): works on a non-paid account

## License

MIT

