# Handle Help

Check if your new handle is available on websites

## Running

```bash
$ go run handlehelp.go
```

```bash
$ go test
```

## Add a New Site

All sites implement the `Website` interface:

```go
type Website struct {
   func checkHandle(handle string) bool, error {}
}
```

Add your site to list of websites in `websites.go`. Make sure you add some
tests to `websites_test.go` as well.
