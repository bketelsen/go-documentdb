## go-documentdb
go-documentdb is a golang implementation of the DocumentDB rest api.  The current implementation is not complete and only provides basic functionality.

## Example Usage

```go
url := "https://xx.documents.azure.com:443"
key := "02y0rthMgYRplBl2ztiRyXQuBFYkXluNDpKf/lNaSJiMKL6AYzwyxjRwdNEFWvvWo4TkpA6i3+T5f8FQEeDf8Q=="
client := NewClient(url, Config{key})
db, err := client.CreateDatabase("foo")
```

## Install

```
go get github.com/nerdylikeme/go-documentdb
```
