package constants

import "encoding/xml"

//VirtualNetworkDefinition - provides the structure of the templated YAML file
type VirtualNetworkDefinition struct {
	Deployment struct {
		DeploymentName string `yaml:"name"`
	} `yaml:"deployment"`
	Networks []NetworkDefinition `yaml:"networks"`
	Host     []HostDefintion     `yaml:"hosts"`
}

//NetworkDefinition - Defines the networks to be built
type NetworkDefinition struct {
	NetworkName string `yaml:"name"`
	NetworkAddr string `yaml:"netaddr"`
	DHCPLower   string `yaml:"dhcplower"`
	DHCPUpper   string `yaml:"dhcpupper"`
	Netmask     string `yaml:"netmask"`
	Type        string `yaml:"type"`
}

// HostDefintion - Defines the host on the virtual network
type HostDefintion struct {
	HostName string   `yaml:"name"`
	Image    string   `yaml:"image"`
	RAM      int      `yaml:"ram"`
	CPUs     int      `yaml:"cpus"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Networks []string `yaml:"networks"`
	HDSpace  string   `yaml:"hd"`
}

//Domain writes the XML files
type Domain struct {
	XMLName xml.Name `xml:"domain"`
	Text    string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name"`
	Memory  struct {
		Value int    `xml:",chardata"`
		Unit  string `xml:"unit,attr"`
	} `xml:"memory"`
	Vcpu struct {
		CPUs      int    `xml:",chardata"`
		Placement string `xml:"placement,attr"`
	} `xml:"vcpu"`
	Os struct {
		Text string `xml:",chardata"`
		Type struct {
			Text string `xml:",chardata"`
			Arch string `xml:"arch,attr"`
		} `xml:"type"`
		Boot struct {
			Text string `xml:",chardata"`
			Dev  string `xml:"dev,attr"`
		} `xml:"boot"`
	} `xml:"os"`
	OnPoweroff string `xml:"on_poweroff"`
	OnReboot   string `xml:"on_reboot"`
	OnCrash    string `xml:"on_crash"`
	Devices    struct {
		Text      string      `xml:",chardata"`
		Emulator  string      `xml:"emulator"`
		Disk      []Disk      `xml:"disk"`
		Interface []Interface `xml:"interface"`
		Serial    struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Target struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
				Port int    `xml:"port,attr"`
			} `xml:"target"`
		} `xml:"serial"`
		Console struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Target struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
				Port int    `xml:"port,attr"`
			} `xml:"target"`
		} `xml:"console"`
	} `xml:"devices"`
}

type Disk struct {
	Text   string `xml:",chardata"`
	Type   string `xml:"type,attr"`
	Device string `xml:"device,attr"`
	Driver struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
		Type string `xml:"type,attr"`
	} `xml:"driver"`
	Source struct {
		Text string `xml:",chardata"`
		File string `xml:"file,attr"`
	} `xml:"source"`
	Target struct {
		Text string `xml:",chardata"`
		Dev  string `xml:"dev,attr"`
		Bus  string `xml:"bus,attr"`
	} `xml:"target"`
}

type Interface struct {
	Text   string `xml:",chardata"`
	Type   string `xml:"type,attr"`
	Source struct {
		Text    string `xml:",chardata"`
		Network string `xml:"network,attr"`
	} `xml:"source"`
	Model struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
		Name string `xml:"name,attr"`
	} `xml:"model"`
}

type Network struct {
	XMLName xml.Name `xml:"network"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name"`
	Forward struct {
		Text string `xml:",chardata"`
		Mode string `xml:"mode,attr"`
	} `xml:"forward"`
	Bridge struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		Stp   string `xml:"stp,attr"`
		Delay string `xml:"delay,attr"`
	} `xml:"bridge"`
	// Mac struct {
	// 	Text    string `xml:",chardata"`
	// 	Address string `xml:"address,attr"`
	// } `xml:"mac"`
	IP struct {
		Text    string `xml:",chardata"`
		Address string `xml:"address,attr"`
		Netmask string `xml:"netmask,attr"`
		Dhcp    struct {
			Text  string `xml:",chardata"`
			Range struct {
				Text  string `xml:",chardata"`
				Start string `xml:"start,attr"`
				End   string `xml:"end,attr"`
			} `xml:"range"`
		} `xml:"dhcp"`
	} `xml:"ip"`
}
