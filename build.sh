#!/bin/bash
echo "Building provider vsphere"
go build -o terraform-provider-vsphere
	
if [ $? -ne 0 ]; then
	echo ""
    echo "ERROR: build failed"
fi
