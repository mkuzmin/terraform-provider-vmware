package vsphere

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"github.com/vmware/govmomi/vim25/types"
)

func TestAccVirtualMachine_basic(t *testing.T) {
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_basic,
			Check:  testAccCheckVirtualMachineState(&vm),
		}},
	},
	)
}

func testAccCheckVirtualMachineState(vm *driver.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["vmware_virtual_machine.vm"]
		if !ok {
			return fmt.Errorf("Not found: %s", "vmware_virtual_machine.vm")
		}

		p := rs.Primary
		if p.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		d, err := driver.NewDriver(
			&driver.ConnectConfig{
				VCenterServer:      "vcenter.vsphere55.test",
				Username:           "root",
				Password:           "jetbrains",
				InsecureConnection: true,
			},
		)
		if err != nil {
			return fmt.Errorf("Cannot connect: %s", err)
		}

		v := d.NewVM(&types.ManagedObjectReference{Type: "VirtualMachine", Value: p.ID})
		*vm = *v

		return nil
	}
}

const testAccVirtualMachine_basic = `
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "empty"
  power_on = false
}
`

func TestAccVirtualMachine_linkedClone(t *testing.T) {
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_linkedClone,
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccCheckVirtualMachineState(&vm),
				testAccCheckLinkedClone(&vm),
			),

		}},
	},
	)
}

func testAccCheckLinkedClone(vm *driver.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vmInfo, err := vm.Info("layoutEx.disk")
		if err != nil {
			return fmt.Errorf("Cannot read VM properties: %v", err)
		}

		if len(vmInfo.LayoutEx.Disk[0].Chain) != 2 {
			return fmt.Errorf("Not a linked clone")
		}

		return nil
	}
}

const testAccVirtualMachine_linkedClone = `
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "basic"
  host = "esxi-1.vsphere55.test"
  linked_clone = true
  power_on = false
}
`
