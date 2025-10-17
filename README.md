# CLI Bot

This is just a test to have a custom CLI tool in the Go language. In a platform engineering environment a tool like this could be used to provide:
- release manager with a simple way to fire releases
- on-call teams to execute same automated checks
- etc

Two systems have been introduced (AWS and Github) with just a mock response.

## Development 

go run .  
go mod tidy # cleans up your module’s dependencies so go.mod and go.sum match what your code actually imports.
go run . aws --service rds

## Build

go build
go build -o hello


## Run

./hello -name Ivan


## Format code

gofmt -w .

./main -help