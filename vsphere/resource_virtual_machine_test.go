package vsphere

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVirtualMachine_basic(t *testing.T) {
	var vm_id string
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_empty,
			Check:  testAccCheckVirtualMachineState(&vm_id),
		}},
	},
	)
}

func testAccCheckVirtualMachineState(vm_id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["vmware_virtual_machine.vm"]
		if !ok {
			return fmt.Errorf("Not found: %s", "vmware_virtual_machine.vm")
		}

		p := rs.Primary
		if p.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		*vm_id = p.ID

		return nil
	}
}

const testAccVirtualMachine_empty = `
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "empty"
  power_on = false
}
`
