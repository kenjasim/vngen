## QEMU Command to start a system

```s
/usr/bin/qemu-system-x86_64 \
-name guest=master1,debug-threads=on \
-S \
-object secret,id=masterKey0,format=raw,file=/var/lib/libvirt/qemu/domain-1-master1/master-key.aes \
-machine pc-i440fx-groovy,accel=kvm,usb=off,dump-guest-core=off \
-cpu qemu64 \
-m 1954 \
-overcommit mem-lock=off \
-smp 2,sockets=2,cores=1,threads=1 \
-uuid b5e27526-692f-4c81-93b4-b854f21a2d1d \
-display none \
-no-user-config \
-nodefaults \
-chardev socket,id=charmonitor,fd=29,server,nowait \
-mon chardev=charmonitor,id=monitor,mode=control \
-rtc base=utc \
-no-shutdown \
-no-acpi \
-boot strict=on \
-device piix3-usb-uhci,id=usb,bus=pci.0,addr=0x1.0x2 \
-blockdev '{"driver":"file","filename":"/var/lib/nenvn/images/ubuntu.img","node-name":"libvirt-3-storage","auto-read-only":true,"discard":"unmap"}' \
-blockdev '{"node-name":"libvirt-3-format","read-only":true,"driver":"qcow2","file":"libvirt-3-storage","backing":null}' \
-blockdev '{"driver":"file","filename":"/var/lib/nenvn/machines/master1/master1.qcow2","node-name":"libvirt-2-storage","auto-read-only":true,"discard":"unmap"}' \
-blockdev '{"node-name":"libvirt-2-format","read-only":false,"driver":"qcow2","file":"libvirt-2-storage","backing":"libvirt-3-format"}' \
-device virtio-blk-pci,bus=pci.0,addr=0x3,drive=libvirt-2-format,id=virtio-disk0,bootindex=1 \
-blockdev '{"driver":"file","filename":"/var/lib/nenvn/machines/master1/master1-seed.qcow2","node-name":"libvirt-1-storage","auto-read-only":true,"discard":"unmap"}' \
-blockdev '{"node-name":"libvirt-1-format","read-only":false,"driver":"raw","file":"libvirt-1-storage"}' \
-device virtio-blk-pci,bus=pci.0,addr=0x4,drive=libvirt-1-format,id=virtio-disk1 \
-netdev tap,fd=31,id=hostnet0,vhost=on,vhostfd=32 \
-device virtio-net-pci,netdev=hostnet0,id=net0,mac=52:54:00:06:69:f0,bus=pci.0,addr=0x2 \
-chardev pty,id=charserial0 \
-device isa-serial,chardev=charserial0,id=serial0 \
-device virtio-balloon-pci,id=balloon0,bus=pci.0,addr=0x5 \
-sandbox on,obsolete=deny,elevateprivileges=deny,spawn=deny,resourcecontrol=deny \
-msg timestamp=on
```


```s
sudo qemu-system-aarch64 \
-smp 2 \
-m 1024 \
-M virt \
-cpu cortex-a57 \
-nographic \
-blockdev '{"driver":"file","filename":"/var/lib/nenvn/images/ubuntu.img","node-name":"libvirt-3-storage","auto-read-only":true,"discard":"unmap"}' \
-blockdev '{"node-name":"libvirt-3-format","read-only":true,"driver":"qcow2","file":"libvirt-3-storage","backing":null}' \
-blockdev '{"driver":"file","filename":"/var/lib/nenvn/machines/master1/master1.qcow2","node-name":"libvirt-2-storage","auto-read-only":true,"discard":"unmap"}' \
-blockdev '{"node-name":"libvirt-2-format","read-only":false,"driver":"qcow2","file":"libvirt-2-storage","backing":"libvirt-3-format"}' \
-device virtio-blk-pci,addr=0x3,drive=libvirt-2-format,id=virtio-disk0,bootindex=1 \
-blockdev '{"driver":"file","filename":"/var/lib/nenvn/machines/master1/master1-seed.qcow2","node-name":"libvirt-1-storage","auto-read-only":true,"discard":"unmap"}' \
-blockdev '{"node-name":"libvirt-1-format","read-only":false,"driver":"raw","file":"libvirt-1-storage"}' \
-device virtio-blk-pci,addr=0x4,drive=libvirt-1-format,id=virtio-disk1 \
-device virtio-balloon-pci,id=balloon0,addr=0x5 
```
