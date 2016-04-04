#!/bin/bash
set -ev
go test ./...
echo ${TRAVIS_GO_VERSION}
if [ "${TRAVIS_GO_VERSION}" = "go1.6" ]; then
 GOARCH=386 go test ./...
 go tool vet -example .
 go get github.com/client9/misspell/cmd/misspell
 go get github.com/soniakeys/vetc
 misspell * **/*
 vetc
fi
