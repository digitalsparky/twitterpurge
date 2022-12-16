### Twitter Purge

This program deletes all tweets from a users twitter account.

#### Requirements

- Golang 1.19.4+ (if building from source)
- Twitter API credentials (see below)

#### Configuration

*** Get API Credentials ***
- Go to https://developer.twitter.com/
- Go to 'developer portal'
- Add app
- Choose 'production'
- Fill in details like 'account maintenance' or something like that
- Fetch your API and consumer keys

*** Copy & update config ***
Copy .env_example file to .env and update with the relevant API details from the previous step.

### Build from source

#### Fetch the dependencies
```
go mod tidy
```

#### Run the purge (you may need to run this multiple times)
```
go build cmd/purge.go
```