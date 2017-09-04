#!/bin/sh

set -e

rm -f bin/*
rm -rf vendor/*

glide install -v

export CGO_ENABLED=0
export GOARCH=amd64

mkdir -p bin
GOOS=darwin  go build -o bin/terraform-provider-vmware.macos
GOOS=linux   go build -o bin/terraform-provider-vmware.linux
GOOS=windows go build -o bin/terraform-provider-vmware.exe

go run schema/schema.go vmware bin
