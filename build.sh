#!/bin/sh -eux

glide install -v

GOOS=darwin  GOARCH=amd64 go build -o bin/terraform-provider-vmware_macos_x64
GOOS=linux   GOARCH=amd64 go build -o bin/terraform-provider-vmware_linux_x64
GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-vmware.exe_windows_x64
