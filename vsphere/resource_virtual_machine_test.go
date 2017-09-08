package vsphere

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jetbrains-infra/packer-builder-vsphere/driver"
	"github.com/vmware/govmomi/vim25/types"
	"os"
)

func TestAccVirtualMachine_basic(t *testing.T) {
	var d driver.Driver
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_basic,
			Check:  testAccCheckVirtualMachineState(&d, &vm),
		}},
	},
	)
}

func testAccCheckVirtualMachineState(driver *driver.Driver, vm *driver.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["vmware_virtual_machine.test"]
		if !ok {
			return fmt.Errorf("Not found: %s", "vmware_virtual_machine.test")
		}

		p := rs.Primary
		if p.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		d, err := newDriver()
		if err != nil {
			return fmt.Errorf("Cannot connect: %s", err)
		}
		*driver = *d

		v := d.NewVM(&types.ManagedObjectReference{Type: "VirtualMachine", Value: p.ID})
		*vm = *v

		return nil
	}
}

const testAccVirtualMachine_basic = `
resource "vmware_virtual_machine" "test" {
  name =  "vm-1"
  image = "empty"
  power_on = false
}
`

func TestAccVirtualMachine_IP(t *testing.T) {
	var d driver.Driver
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_IP,
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccCheckVirtualMachineState(&d, &vm),
				testAccCheckIP(&vm),
			),
		}},
	},
	)
}

func testAccCheckIP(vm *driver.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vmInfo, err := vm.Info("guest.ipAddress")
		if err != nil {
			return fmt.Errorf("Cannot read VM properties: %v", err)
		}

		v, ok := getAttribute(s, "ip_address")
		if !ok {
			return fmt.Errorf("Attribute 'ip_address' not found")
		}

		if vmInfo.Guest.IpAddress != v {
			return fmt.Errorf("invalid IP address")
		}

		return nil
	}
}

const testAccVirtualMachine_IP = `
resource "vmware_virtual_machine" "test" {
  name =  "vm-1"
  image = "basic"
  host = "esxi-1.vsphere55.test"
  linked_clone = true
}
`

func TestAccVirtualMachine_linkedClone(t *testing.T) {
	var d driver.Driver
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_linkedClone,
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccCheckVirtualMachineState(&d, &vm),
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
resource "vmware_virtual_machine" "test" {
  name =  "vm-1"
  image = "basic"
  host = "esxi-1.vsphere55.test"
  linked_clone = true
  power_on = false
}
`
func TestAccVirtualMachine_pool(t *testing.T) {
	var d driver.Driver
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_pool,
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccCheckVirtualMachineState(&d, &vm),
				testAccCheckPool(&d, &vm),
			),
		}},
	},
	)
}

func testAccCheckPool(d *driver.Driver, vm *driver.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vmInfo, err := vm.Info("resourcePool")
		if err != nil {
			return fmt.Errorf("Cannot read VM properties: %v", err)
		}

		p := d.NewResourcePool(vmInfo.ResourcePool)
		path, err := p.Path()
		if err != nil {
			return fmt.Errorf("Cannot read resource pool name: %v", err)
		}

		pool, ok := getAttribute(s, "resource_pool")
		if !ok {
			return fmt.Errorf("Cannot read 'resource_pool' attribute")
		}

		if path != pool {
			return fmt.Errorf("Wrong folder. expected: %v, got: %v", pool, path)
		}

		return nil
	}
}

const testAccVirtualMachine_pool = `
resource "vmware_virtual_machine" "test" {
  name =  "vm-1"
  image = "basic"
  host = "esxi-1.vsphere55.test"
  resource_pool = "pool1/pool2"
  linked_clone = true
  power_on = false
}
`

func TestAccVirtualMachine_folder(t *testing.T) {
	var d driver.Driver
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_folder,
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccCheckVirtualMachineState(&d, &vm),
				testAccCheckFolder(&d, &vm),
			),
		}},
	},
	)
}

func testAccCheckFolder(d *driver.Driver, vm *driver.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vmInfo, err := vm.Info("parent")
		if err != nil {
			return fmt.Errorf("Cannot read VM properties: %v", err)
		}

		f := d.NewFolder(vmInfo.Parent)
		path, err := f.Path()
		if err != nil {
			return fmt.Errorf("Cannot read folder name: %v", err)
		}

		folder, ok := getAttribute(s, "folder")
		if !ok {
			return fmt.Errorf("Cannot read 'folder' attribute")
		}

		if path != folder {
			return fmt.Errorf("Wrong folder. expected: %v, got: %v", folder, path)
		}

		return nil
	}
}

const testAccVirtualMachine_folder = `
resource "vmware_virtual_machine" "test" {
  name =  "vm-1"
  image = "basic"
  host = "esxi-1.vsphere55.test"
  folder = "folder1/folder2"
  linked_clone = true
  power_on = false
}
`

func TestAccVirtualMachine_hardware(t *testing.T) {
	var d driver.Driver
	var vm driver.VirtualMachine
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: testAccVirtualMachine_hardware,
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccCheckVirtualMachineState(&d, &vm),
				testAccCheckHardware(&d, &vm),
			),
		}},
	},
	)
}

func testAccCheckHardware(d *driver.Driver, vm *driver.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vmInfo, err := vm.Info("config.hardware")
		if err != nil {
			return fmt.Errorf("Cannot read VM properties: %v", err)
		}

		if vmInfo.Config.Hardware.NumCPU != 2 {
			return fmt.Errorf("wrong number of CPU sockets")
		}

		if vmInfo.Config.Hardware.MemoryMB != 2048 {
			return fmt.Errorf("wrong RAM size")
		}

		return nil
	}
}

const testAccVirtualMachine_hardware = `
resource "vmware_virtual_machine" "test" {
  name =  "vm-1"
  image = "basic"
  host = "esxi-1.vsphere55.test"

  cpus = 2
  memory = 2048

  linked_clone = true
  power_on = false
}
`

func primaryInstanceState(s *terraform.State, name string) (*terraform.InstanceState, error) {
	ms := s.RootModule()
	rs, ok := ms.Resources[name]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", name)
	}

	is := rs.Primary
	if is == nil {
		return nil, fmt.Errorf("No primary instance: %s", name)
	}

	return is, nil
}

func newDriver() (*driver.Driver, error) {
	d, err := driver.NewDriver(
		&driver.ConnectConfig{
			VCenterServer:      os.Getenv("VSPHERE_SERVER"),
			Username:           os.Getenv("VSPHERE_USER"),
			Password:           os.Getenv("VSPHERE_PASSWORD"),
			InsecureConnection: os.Getenv("VSPHERE_INSECURE") == "true",
		},
	)
	return d, err
}

func getAttribute(s *terraform.State, attr string) (string, bool) {
	is, err := primaryInstanceState(s, "vmware_virtual_machine.test")
	if err != nil {
		return "", false
	}

	v, ok := is.Attributes[attr]
	return v, ok
}