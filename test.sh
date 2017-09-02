#!/bin/sh

set -e

export TF_ACC=1
export VSPHERE_SERVER=vcenter.vsphere55.test
export VSPHERE_USER=root
export VSPHERE_PASSWORD=jetbrains
export VSPHERE_INSECURE=true

go test -v $(go list ./... | grep -v '/vendor/')
