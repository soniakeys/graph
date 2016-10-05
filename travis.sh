#!/bin/bash
set -ex
go test ./...
if [ "$TRAVIS_GO_VERSION" = "1.7" ]; then
 GOARCH=386 go test ./...
 go tool vet -composites=false -printf=false -shift=false .
 go get github.com/client9/misspell/cmd/misspell
 go get github.com/soniakeys/vetc
 misspell -error * */* */*/*
 vetc
fi
