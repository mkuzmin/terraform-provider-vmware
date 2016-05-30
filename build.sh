#!/bin/sh -eux

VERSION=0.6-dev

rm -rf bin
GOOS=darwin  GOARCH=amd64 go build -o bin/macos/terraform-provider-vsphere
GOOS=linux   GOARCH=amd64 go build -o bin/linux/terraform-provider-vsphere
GOOS=windows GOARCH=amd64 go build -o bin/windows/terraform-provider-vsphere.exe

tar czf bin/terraform-vsphere-$VERSION-macos.tar.gz  --directory=bin/macos terraform-provider-vsphere
tar czf bin/terraform-vsphere-$VERSION-linux.tar.gz  --directory=bin/linux terraform-provider-vsphere
zip     bin/terraform-vsphere-$VERSION-windows.zip   -j bin/windows/terraform-provider-vsphere.exe
