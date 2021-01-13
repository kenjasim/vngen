package host

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"nenvoy.com/pkg/database"
	"nenvoy.com/pkg/utils/printing"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	libvirt "libvirt.org/libvirt-go"
	"nenvoy.com/pkg/constants"
	structs "nenvoy.com/pkg/constants"
	cmd "nenvoy.com/pkg/utils/cmd"
	"nenvoy.com/pkg/utils/files"
)

var errNameUsed = errors.New("Host name already used")

//Host - Struct for the host data in the database
type Host struct {
	gorm.Model
	Name         string
	Image        string
	RAM          int
	CPUs         int
	Username     string
	Password     string
	HDSpace      string
	DeploymentID uint
}

// createHostXML - Create the host domain
func (h *Host) createHostXML(networks []string) (domainDef string, err error) {
	//Define the domain object for libvirt
	domain := structs.Domain{}

	// Set the metadata values
	domain.Type = "kvm"
	domain.Name = h.Name
	// Memory values
	domain.Memory.Unit = "MB"
	domain.Memory.Value = h.RAM
	// CPU Values
	domain.Vcpu.Placement = "static"
	domain.Vcpu.CPUs = h.CPUs
	// Set the OS type (hvm means hardware virtualisation)
	domain.Os.Type.Arch = "x86_64"
	domain.Os.Type.Text = "hvm"
	domain.Os.Boot.Dev = "hd"

	// Set the actions to be completed on a VM action
	domain.OnPoweroff = "destroy"
	domain.OnReboot = "restart"
	domain.OnCrash = "restart"

	// Set the device values
	domain.Devices.Emulator = "/usr/bin/qemu-system-x86_64"

	// Create the disks that are required
	err = h.createHostDisks()
	if err != nil {
		return "", err
	}

	// Main hard drive with all the content
	mainHD := structs.Disk{}
	mainHD.Type = "file"
	mainHD.Device = "disk"
	mainHD.Driver.Name = "qemu"
	mainHD.Driver.Type = "qcow2"
	mainHD.Source.File = fmt.Sprintf("/var/lib/nenvn/machines/%s/%s.qcow2", h.Name, h.Name)
	mainHD.Target.Dev = "vda"
	mainHD.Target.Bus = "virtio"

	domain.Devices.Disk = append(domain.Devices.Disk, mainHD)

	// Cloud Init hard drive with all the user defined info
	cloudInitHD := structs.Disk{}
	cloudInitHD.Type = "file"
	cloudInitHD.Device = "disk"
	cloudInitHD.Driver.Name = "qemu"
	cloudInitHD.Driver.Type = "raw"
	cloudInitHD.Source.File = fmt.Sprintf("/var/lib/nenvn/machines/%s/%s-seed.qcow2", h.Name, h.Name)
	cloudInitHD.Target.Dev = "vdb"
	cloudInitHD.Target.Bus = "virtio"
	domain.Devices.Disk = append(domain.Devices.Disk, cloudInitHD)

	// Setup the interfaces
	for _, network := range networks {
		iface := structs.Interface{}
		iface.Type = "network"
		iface.Source.Network = network
		iface.Model.Name = "isa_serial"
		iface.Model.Type = "virtio"
		domain.Devices.Interface = append(domain.Devices.Interface, iface)
	}

	//Serial and console connection
	domain.Devices.Serial.Type = "pty"
	domain.Devices.Serial.Target.Type = "isa-serial"
	domain.Devices.Serial.Target.Port = 0
	domain.Devices.Console.Type = "pty"
	domain.Devices.Console.Target.Type = "serial"
	domain.Devices.Console.Target.Port = 0

	xmlBytes, err := xml.MarshalIndent(domain, "", "	")
	if err != nil {
		return "", err
	}

	return string(xmlBytes), nil
}

// Start - starts the VM
func (h *Host) Start() (err error) {
	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get the domain
	dom, err := conn.LookupDomainByName(h.Name)
	if err != nil {
		return err
	}

	// Start the domain
	err = dom.Create()
	if err != nil {
		return err
	}

	printing.PrintSuccess(fmt.Sprintf("Started host %s", h.Name))
	return nil
}

// Restart - restarts the VM
func (h *Host) Restart() (err error) {
	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get the domain
	dom, err := conn.LookupDomainByName(h.Name)
	if err != nil {
		return err
	}

	// Start the domain
	err = dom.Reboot(1)
	if err != nil {
		return err
	}

	printing.PrintSuccess(fmt.Sprintf("Restarted host %s", h.Name))
	return nil
}

// Stop - stops the VM
func (h *Host) Stop() (err error) {
	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get the domain
	dom, err := conn.LookupDomainByName(h.Name)
	if err != nil {
		return err
	}

	// Start the domain
	err = dom.Destroy()
	if err != nil {
		return err
	}

	printing.PrintSuccess(fmt.Sprintf("Stopped host %s", h.Name))
	return nil
}

// Destroy - destroy the VM
func (h *Host) Destroy() (err error) {
	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get the domain
	dom, err := conn.LookupDomainByName(h.Name)
	if err != nil {
		return err
	}

	// Get the domain status
	domState, _, err := dom.GetState()
	if err != nil {
		return err
	}

	// If the domain is running stop it
	if domState == 1 {
		err = dom.Destroy()
		if err != nil {
			return err
		}
	}

	// Undefine the domain
	err = dom.Undefine()
	if err != nil {
		return err
	}

	db, err := database.NewSession()
	if err != nil {
		return err
	}

	// Remove from the database
	db.Delete(&h)

	// Remove the machine directory
	err = files.RemoveDirectories([]string{fmt.Sprintf("%s/machines/%s", constants.AppDir, h.Name)})
	if err != nil {
		return errors.Wrap(err, "failed to remove directories")
	}

	printing.PrintSuccess(fmt.Sprintf("Destroyed host %s", h.Name))
	return nil
}

// CreateHost - Creates the host domain from the XML template
func (h *Host) CreateHost(networks []string) (err error) {
	// Get the xml hosts
	hostDef, err := h.createHostXML(networks)
	if err != nil {
		return err
	}

	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Define the domain
	_, err = conn.DomainDefineXML(hostDef)
	if err != nil {
		return err
	}
	return nil
}

// GetHostState - returns the VMState
func (h *Host) GetHostState() (state string, err error) {

	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return state, err
	}
	defer conn.Close()

	// Get the domain by name
	dom, err := conn.LookupDomainByName(h.Name)
	if err != nil {
		return state, err
	}
	defer dom.Free()

	// REF: Domain states https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainState
	domState, _, err := dom.GetState()
	if domState == 5 {
		state = "off"
	} else if domState == 1 {
		state = "running"
	}

	return state, nil
}

// GetHostIfaces - returns the Host IPs
func (h *Host) GetHostIfaces() (ifaces map[string][]string, err error) {
	// Connect to the libvirt socket
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return ifaces, err
	}
	defer conn.Close()

	ifaces = make(map[string][]string)

	// Get the domain by name
	dom, err := conn.LookupDomainByName(h.Name)
	if err != nil {
		return ifaces, err
	}
	defer dom.Free()

	// REF: Domain states https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainState
	domState, _, err := dom.GetState()
	if domState != 1 {
		return
	}

	// Get the interfaces
	domIfaces, err := dom.ListAllInterfaceAddresses(0)

	for _, iface := range domIfaces {
		// Loop through the interface and assign the IPs
		macIP := []string{iface.Hwaddr}
		for _, ip := range iface.Addrs {
			macIP = append(macIP, ip.Addr)
		}

		// Add to map
		ifaces[iface.Name] = macIP
	}

	return ifaces, nil
}

// CreateHostDisks - Create the host disks which is needed for the vm
func (h *Host) createHostDisks() (err error) {
	// Create host directory
	dirs := []string{fmt.Sprintf("/var/lib/nenvn/machines/%s", h.Name)}
	err = files.CreateDirectories(dirs)
	if err != nil {
		return err
	}

	// Create the VM main image
	_, stderr, err := cmd.Output("qemu-img", "create", "-F", "qcow2", "-b", fmt.Sprintf("/var/lib/nenvn/images/%s.img", h.Image), "-f", "qcow2", fmt.Sprintf("/var/lib/nenvn/machines/%s/%s.qcow2", h.Name, h.Name), h.HDSpace)
	if err != nil {
		return errors.Wrap(err, stderr)
	}

	// Create the user-data File
	userData := `#cloud-config
hostname: %s
manage_etc_hosts: true
users:
  - name: %s
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: users, admin
    home: /home/%s
    shell: /bin/bash
    lock_passwd: false
ssh_pwauth: true
disable_root: false
chpasswd:
  list: |
    %s:%s
  expire: false
`

	//Add the variables and write the files
	userData = fmt.Sprintf(userData, h.Name, h.Username, h.Username, h.Username, h.Password)
	err = ioutil.WriteFile(fmt.Sprintf("/var/lib/nenvn/machines/%s/user-data", h.Name), []byte(userData), 0755)

	// Create an empty meta-data file
	err = ioutil.WriteFile(fmt.Sprintf("/var/lib/nenvn/machines/%s/meta-data", h.Name), []byte(""), 0755)

	// Create the cloud-init disk
	_, stderr, err = cmd.Output("cloud-localds", "-v", fmt.Sprintf("/var/lib/nenvn/machines/%s/%s-seed.qcow2", h.Name, h.Name), fmt.Sprintf("/var/lib/nenvn/machines/%s/user-data", h.Name), fmt.Sprintf("/var/lib/nenvn/machines/%s/meta-data", h.Name))
	if err != nil {
		return errors.Wrap(err, stderr)
	}

	return nil
}

// DefineHost - defines the host and writes the XML config file
func DefineHost(hostDef structs.HostDefintion) (host Host, err error) {

	// Check if the name exists in the database
	hostTest, err := GetHostByName(hostDef.HostName)
	if hostTest != (Host{}) {
		return host, errNameUsed
	}

	// Create host struct for database
	host = Host{
		Name:     hostDef.HostName,
		Image:    hostDef.Image,
		RAM:      hostDef.RAM,
		CPUs:     hostDef.CPUs,
		Username: hostDef.Username,
		Password: hostDef.Password,
		HDSpace:  hostDef.HDSpace,
	}

	return host, nil
}

// GetHosts - returns all the hosts in the database
func GetHosts() (hosts []Host, err error) {
	// Connect and open the database
	db, err := database.NewSession()
	if err != nil {
		return nil, err
	}

	err = db.Find(&hosts).Error
	if err != nil {
		return hosts, errors.Wrap(err, "could not find hosts")
	}

	return hosts, nil
}

//GetHostsByDeployment - returns all the hosts in a deployment
func GetHostsByDeployment(ID uint) (hosts []Host, err error) {
	// Connect and open the database
	db, err := database.NewSession()
	if err != nil {
		return nil, err
	}

	err = db.Where("deployment_id = ?", ID).Find(&hosts).Error
	if err != nil {
		return hosts, errors.Wrap(err, "could not find hosts")
	}

	return hosts, nil
}

// GetHostByName - Returns the host with that particular name
func GetHostByName(name string) (host Host, err error) {
	// Connect and open the database
	db, err := database.NewSession()
	if err != nil {
		return host, err
	}

	err = db.Where("name = ?", name).First(&host).Error
	if err == gorm.ErrRecordNotFound {
		return host, nil
	} else if err != nil {
		return host, errors.Wrap(err, "could not find host")
	}

	return host, nil
}
