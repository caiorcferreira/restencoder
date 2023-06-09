# restencoder
[![Go Reference](https://pkg.go.dev/badge/github.com/caiorcferreira/restencoder.svg)](https://pkg.go.dev/github.com/caiorcferreira/restencoder)
[![Go Report Card](https://goreportcard.com/badge/github.com/caiorcferreira/restencoder)](https://goreportcard.com/report/github.com/caiorcferreira/restencoder)
[![Builds](https://github.com/caiorcferreira/restencoder/actions/workflows/main.yaml/badge.svg)](https://github.com/caiorcferreira/restencoder/actions/workflows/main.yaml)

A optioned REST response encoder compatible with Go's standard lib. 

## Installation
```
go get github.com/caiorcferreira/restencoder
```

## Usage

### Success response
```go
restencoder.Respond(
    w,
    StatusCode(http.StatusCreated),
    Header("X-Header", "value"),
    JSONBody(body),
)
```

### Failure response
```go
restencoder.Respond(
    w,
    StatusCode(http.StatusBadRequest),
    Header("X-Header", "value"),
    Error(err),
)
```

If you want a different failure response schema, you could implement your own `restencoder.ResponseOption` to define `ResponseConfig.JSONBody` to your failure response struct.

## Contributing
Every help is always welcome. Feel free do throw us a pull request, we'll do our best to check it out as soon as possible. But before that, let us establish some guidelines:

1. This is an open source project so please do not add any proprietary code or infringe any copyright of any sort.
2. Avoid unnecessary dependencies or messing up go.mod file.
3. Be aware of golang coding style. Use a lint to help you out.
4.  Add tests to cover your contribution.
5. Use meaningful [messages](https://medium.com/@menuka/writing-meaningful-git-commit-messages-a62756b65c81) to your commits.
6. Use [pull requests](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests).
7. At last, but also important, be kind and polite with the community.

Any submitted issue which disrespect one or more guidelines above, will be discarded and closed.


## License

Released under the [MIT License](LICENSE).