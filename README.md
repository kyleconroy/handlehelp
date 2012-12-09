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

Each sites has a function that matches this signature

```go
func website(handle string) (string, bool) {
  return "example.com", false
}
```

Add your site function to list of websites in `handlehelp.go`.
