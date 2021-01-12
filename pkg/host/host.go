package host

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	libvirt "libvirt.org/libvirt-go"
	"nenvoy.com/pkg/constants"
	structs "nenvoy.com/pkg/constants"
	cmd "nenvoy.com/pkg/utils/cmd"
	"nenvoy.com/pkg/utils/files"
)

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

// DefineHost - defines the host and writes the XML config file
func DefineHost(hostDef structs.HostDefintion) (host Host, domainDef string, err error) {

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

	domainDef, err = createHostXML(host, hostDef.Networks)
	if err != nil {
		return host, domainDef, err
	}

	return host, domainDef, nil
}

// createHostXML - Create the host domain
func createHostXML(host Host, networks []string) (domainDef string, err error) {
	//Define the domain object for libvirt
	domain := structs.Domain{}

	// Set the metadata values
	domain.Type = "kvm"
	domain.Name = host.Name
	// Memory values
	domain.Memory.Unit = "MB"
	domain.Memory.Value = host.RAM
	// CPU Values
	domain.Vcpu.Placement = "static"
	domain.Vcpu.CPUs = host.CPUs
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
	err = createHostDisks(host.Name, host.Image, host.HDSpace, host.Username, host.Password)
	if err != nil {
		return "", err
	}

	// Main hard drive with all the content
	mainHD := structs.Disk{}
	mainHD.Type = "file"
	mainHD.Device = "disk"
	mainHD.Driver.Name = "qemu"
	mainHD.Driver.Type = "qcow2"
	mainHD.Source.File = fmt.Sprintf("/var/lib/nenvn/machines/%s/%s.qcow2", host.Name, host.Name)
	mainHD.Target.Dev = "vda"
	mainHD.Target.Bus = "virtio"

	domain.Devices.Disk = append(domain.Devices.Disk, mainHD)

	// Cloud Init hard drive with all the user defined info
	cloudInitHD := structs.Disk{}
	cloudInitHD.Type = "file"
	cloudInitHD.Device = "disk"
	cloudInitHD.Driver.Name = "qemu"
	cloudInitHD.Driver.Type = "raw"
	cloudInitHD.Source.File = fmt.Sprintf("/var/lib/nenvn/machines/%s/%s-seed.qcow2", host.Name, host.Name)
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

// CreateHostDisks - Create the host disks which is needed for the vm
func createHostDisks(name string, image string, space string, username string, password string) (err error) {
	// Create host directory
	dirs := []string{fmt.Sprintf("/var/lib/nenvn/machines/%s", name)}
	err = files.CreateDirectories(dirs)
	if err != nil {
		return err
	}

	// Create the VM main image
	_, stderr, err := cmd.Output("qemu-img", "create", "-F", "qcow2", "-b", fmt.Sprintf("/var/lib/nenvn/images/%s.img", image), "-f", "qcow2", fmt.Sprintf("/var/lib/nenvn/machines/%s/%s.qcow2", name, name), space)
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
	userData = fmt.Sprintf(userData, name, username, username, username, password)
	err = ioutil.WriteFile(fmt.Sprintf("/var/lib/nenvn/machines/%s/user-data", name), []byte(userData), 0755)

	// Create an empty meta-data file
	err = ioutil.WriteFile(fmt.Sprintf("/var/lib/nenvn/machines/%s/meta-data", name), []byte(""), 0755)

	// Create the cloud-init disk
	_, stderr, err = cmd.Output("cloud-localds", "-v", fmt.Sprintf("/var/lib/nenvn/machines/%s/%s-seed.qcow2", name, name), fmt.Sprintf("/var/lib/nenvn/machines/%s/user-data", name), fmt.Sprintf("/var/lib/nenvn/machines/%s/meta-data", name))
	if err != nil {
		return errors.Wrap(err, stderr)
	}

	return nil
}

// CreateHost - Creates the host domain from the XML template
func CreateHost(hostDef string) (err error) {
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

// GetHosts - returns all the hosts in the database
func GetHosts() (hosts []Host, err error) {
	// Connect and open the database
	db, err := gorm.Open(sqlite.Open(constants.DBPath), &gorm.Config{})
	if err != nil {
		return hosts, errors.Wrap(err, "failed to connect database")
	}

	err = db.Find(&hosts).Error
	if err != nil {
		return hosts, errors.Wrap(err, "could not find hosts")
	}

	return hosts, nil
}
