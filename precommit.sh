#!/usr/bin/env sh

echo "gofmt..."
gofmt -w *.go

echo "go lint..."
go get github.com/golang/lint/golint
golint ./...

echo "go vet..."
go get golang.org/x/tools/cmd/vet
go vet ./...

echo "dependencies..."
go list -f '{{join .Deps "\n"}}' | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'

echo "go test..."
go test -coverprofile=/tmp/cover.out ./... && go tool cover -html=/tmp/cover.out -o /tmp/coverage.html && xdg-open /tmp/coverage.html
