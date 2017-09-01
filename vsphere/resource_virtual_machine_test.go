package vsphere

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
	//    "github.com/vmware/govmomi"
	//    "github.com/vmware/govmomi/vim25/types"
)

func TestAccVirtualMachine_basic(t *testing.T) {
	var vm_id string
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVirtualMachine_empty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMachineState(&vm_id),
				),
			},
		},
		//        CheckDestroy: resource.ComposeTestCheckFunc(
		//            testAccCheckVirtualMachineDestroy(&vm_id),
		//        ),
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
		//        if p.Attributes["ip_address"] == "" {
		//            return fmt.Errorf("IP address is not set")
		//        }
		*vm_id = p.ID

		return nil
	}
}

//func testAccCheckVirtualMachineDestroy(vm_id *string) resource.TestCheckFunc {
//    return func(s *terraform.State) error {
//        client := testAccProvider.Meta().(*govmomi.Client)
//
//        vm_mor := types.ManagedObjectReference{Type: "VirtualMachine", Value: *vm_id }
//        err := client.Properties(vm_mor, []string{"summary"}, &vm_mor)
//        if err == nil {
//            return fmt.Errorf("Record still exists")
//        }
//
//        return nil
//    }
//}

const testAccVirtualMachine_empty = `
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "empty"
  power_on = false
}
`
