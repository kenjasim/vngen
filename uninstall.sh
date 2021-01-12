sudo virsh destroy master1
sudo virsh destroy master2
sudo virsh undefine master1 
sudo virsh undefine master2
sudo virsh net-undefine br0
sudo virsh net-destroy br0
sudo rm -r /var/lib/nenvn/machines
sudo rm  /var/lib/nenvn/main.db