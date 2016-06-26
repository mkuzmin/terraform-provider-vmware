#!/bin/sh -eux

VERSION=0.7-dev

rm -rf bin
GOOS=darwin  GOARCH=amd64 go build -o bin/macos/terraform-provider-vmware
GOOS=linux   GOARCH=amd64 go build -o bin/linux/terraform-provider-vmware
GOOS=windows GOARCH=amd64 go build -o bin/windows/terraform-provider-vmware.exe

tar czf bin/terraform-vsphere-$VERSION-macos.tar.gz  --directory=bin/macos terraform-provider-vmware
tar czf bin/terraform-vsphere-$VERSION-linux.tar.gz  --directory=bin/linux terraform-provider-vmware
zip     bin/terraform-vsphere-$VERSION-windows.zip   -j bin/windows/terraform-provider-vmware.exe
