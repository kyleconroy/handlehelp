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

Each sites has an entry in `sites.json` that looks like:

```json
  {
    "Name": "example",
    "UserURL": "http://example.com/user/%s",
    "RegisterURL": "http://example.com/join",
    "Pattern": "^[a-zA-Z0-9_-]{2,15}$",
  },
```

- **Name**: Name of the website
- **UserURL**: URL to check if a user has already taken that name
- **RegisterURL**: URL to sign up for the website
- **Pattern**: Regular expression the handle must match
