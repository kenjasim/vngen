<!-- omit in toc -->
# Virtual Network Generator (vngen)
Automated virtual network generation in golang utilising QEMU virtulisation and KVM hardware acceleration

* Command line based application 
* Automated network virtualisation with KVM and QEMU 
* YAML-defined network-host topology configuration
* Create, run and manage multiple deployment topologies
* Data persistence 
* RestApi server mode for remote connection.
* Log management 
* Supported on Linux only 

<!-- omit in toc -->
## Table of Contents
- [Requirements](#requirements)
  - [Ubuntu/Debian](#ubuntudebian)
  - [Arch Linux](#arch-linux)
  - [Ubuntu/Debian](#ubuntudebian-1)
- [YAML Topology Configuration](#yaml-topology-configuration)
- [Command Line Interface](#command-line-interface)
  - [Installation](#installation)
  - [Create Network Deployment](#create-network-deployment)
  - [Start, Stop, Restart and Destroy Hosts or Deployments](#start-stop-restart-and-destroy-hosts-or-deployments)
  - [Display Information](#display-information)
    - [Hosts](#hosts)
    - [Networks](#networks)
    - [IPs](#ips)
- [Rest API Server](#rest-api-server)
  - [Server (Localhost mode with http)](#server-localhost-mode-with-http)
  - [URL Endpoints](#url-endpoints)
    - [Build](#build)
    - [Start, Stop, Restart, Destroy](#start-stop-restart-destroy)
    - [Details](#details)

## Requirements

Install KVM and Qemu

### Ubuntu/Debian

```shell
sudo apt install -v qemu libvirt-client libvirt-dev cloud-utils
```

### Arch Linux

```shell
sudo pacman -S qemu libvirt cloud-utils
```

### Ubuntu/Debian

## YAML Topology Configuration 

Create custom topologies with yaml templates. Simply add networks, hosts and then assign networks to hosts as per the default.yaml example:

```yaml
---
# Deploymnet Configuration
deployment: 
  name: default

# Network Configuration
networks:
  - name: br0
    netaddr: "20.0.0.1"
    dhcplower: "20.0.0.2"
    dhcpupper: "20.0.0.254"
    netmask: "255.255.255.0"
    type: "nat"

#Hosts Configuration
hosts:
  - name: master1
    image: ubuntu
    ram: 2048
    cpus: 2
    hd: "10G"
    username: dev
    password: ved
    networks:
      - br0
  - name: master2
    image: ubuntu
    ram: 2048
    cpus: 2
    hd: "10G"
    username: dev
    password: ved
    networks:
      - br0
```

## Command Line Interface

### Installation 

```bash
# Clone repository
$ git clone https://github.com/Kenjasim/vngen.git
$ cd vngen 
# Build application binaries
$ go build -o vngen cmd/vngen/vngen.go 
# Add vngen application to PATH
$ mv vngen /usr/local/bin
``` 

### Create Network Deployment
```go
sudo vngen build </path/to/template> # default.yaml
```

### Start, Stop, Restart and Destroy Hosts or Deployments
```go
sudo vngen start [deployment|host] <name>
sudo vngen stop [deployment|host] <name>
sudo vngen restart [deployment|host] <name>
sudo vngen destroy [deployment|host] <name>
```

### Display Information
You can display hosts, networks, and IPs
```go
sudo vngen get hosts
sudo vngen get networks
sudo vngen get ips
```

#### Hosts

```bash
[i] Getting hosts
Name    VMState Image  RAM  CPU Storage Deployment 
master1 off     ubuntu 2048 2   10G     default
master2 off     ubuntu 2048 2   10G     default
```

#### Networks

```bash
[i] Getting networks
Name Type IP       DHCP Range            Deployment 
br0  nat  20.0.0.1 20.0.0.2 - 20.0.0.254 default
```

#### IPs
```
[i] Getting ips
Name    Interface MacAddr           IPs        Deployment 
master1 vnet0     52:54:00:2b:76:8f 20.0.0.141 default
master2 vnet1     52:54:00:ed:90:f9 20.0.0.88  default
```

## Rest API Server 
The RestAPI Server enables remote access to the application either through AVN's client mode, or via direct http (localhost), https (remote) requests. 


### Server (Localhost mode with http) 
To launch the server just run the command on the binary

```
sudo vngen api

# Or you can specify the port

sudo go run cmd/vngen/vngen.go api -p 30675 
```

### URL Endpoints

#### Build

You need to send a `PUT` request with a jsonified version of a YAML template file, an example can be seen below:

```json
{
    "deployment": {
      "name": "default"
    },
    "networks": [
      {
        "name": "br0",
        "netaddr": "20.0.0.1",
        "dhcplower": "20.0.0.2",
        "dhcpupper": "20.0.0.254",
        "netmask": "255.255.255.0",
        "type": "nat"
      }
    ],
    "hosts": [
      {
        "name": "master1",
        "image": "ubuntu",
        "ram": 2048,
        "cpus": 2,
        "hd": "10G",
        "username": "dev",
        "password": "ved",
        "networks": [
          "br0"
        ]
      },
      {
        "name": "master2",
        "image": "ubuntu",
        "ram": 2048,
        "cpus": 2,
        "hd": "10G",
        "username": "dev",
        "password": "ved",
        "networks": [
          "br0"
        ]
      }
    ]
  }
```

You can then pass this to the endpoint as a post request with the header set as `application/json`

```
http://localhost:8000/build
```

#### Start, Stop, Restart, Destroy

These 4 commands follow the same trend, all must be sent as `POST` requests and all have a simmilar structure.

```
http://localhost:8000/<start|stop|restart|destroy>/<deployment|host>/[name]
http://localhost:8000/start/deployment/default
```

#### Details

To get a list of all defined hosts or networks you can use this URL endpoint:

```
http://localhost:8000/hosts
http://localhost:8000/networks
```

To get more details about one host you can use:

```
http://localhost:8000/details/[name]
```

Finally, to get the IP of a host you can use:

```
http://localhost:8000/details/[name]/ipv4
```

These should all be run as `GET` requests


