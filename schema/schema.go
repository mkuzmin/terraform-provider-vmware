package main

import (
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/utils/intellij-hcl/terraform/provider-schema-generator"
    "github.com/mkuzmin/terraform-provider-vmware/vsphere"
)

func main() {
    provider_schema_generator.Generate(vsphere.Provider().(*schema.Provider))
}
