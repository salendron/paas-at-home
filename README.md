# home-service-host
Is it possible to use a RaspberryPi as service host for private home use? Shared Todolists, Automations, Microservices, Docker, all the good stuff...

## The Setup
I'm using a Raspberry Pi 3 Model B and Raspberry Pi OS (32-bit) Lite (Minimal image based on Debian Buster). Simply follow the instructions on [https://www.raspberrypi.org/](https://www.raspberrypi.org/) to setup your pi. 

### Network
After that I've activated SSH and set WiFi using this [guide](https://www.raspberrypi.org/documentation/configuration/wireless/wireless-cli.md). 

#### /etc/wpa_supplicant/wpa_supplicant.conf
```
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
country=AT

network={
        ssid="YOUR-SSID"
        psk="YOUR-PASSWORD"
}
```
#### /etc/network/interfaces
```
# interfaces(5) file used by ifup(8) and ifdown(8)

# Please note that this file is written to be used with dhcpcd
# For static IP, consult /etc/dhcpcd.conf and 'man dhcpcd.conf'

# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

allow-hotplug wlan0
iface wlan0 inet manual
wpa-roam /etc/wpa_supplicant/wpa_supplicant.conf
iface default inet static
        address 192.168.1.82
        netmask 255.255.255.0
        network 192.168.1.1
        gateway 192.168.1.1
```

At last activate wpa_supplicant@.service to start WiFi on boot.

```sudo systemctl enable wpa_supplicant@.service```

### Setup Docker
```
# Update OS
sudo apt-get update
sudo apt-get upgrade

# Install Dependencies
sudo apt-get install -y libffi-dev libssl-dev
sudo apt-get install -y python3 python3-pip
sudo apt-get remove python-configparser

# Install Docker
curl -sSL https://get.docker.com | sh

# Add User pi to group docker
sudo usermod -aG docker pi

# Test Docker
docker run hello-world
```

### Add external storage
First connect your USB drive to the Raspberry Pi.
Let's find it using this command. Also copy the UUId of the device.

```
sudo blkid -o list -w /dev/null
```
Create a mountpoint and then allow user pi to edit the contents of this new directory. (chown, chgrp, ...)
```
sudo mkdir /media/external
```
Maybe you want to format your drive now, but make sure that it is really also sda innyour case!
```
sudo mkfs.ext4 /dev/sda
```
To mount the drive on boot add the following line to /etc/fstab.
```
UUID=YOUR-DEVICE-UUID /media/external/ ext4 defaults 0
```
Now let's test the mountpoint.
```
sudo mount -a
```
Your USB storage should now be mounted and also automatically get mounted after a reboot.
