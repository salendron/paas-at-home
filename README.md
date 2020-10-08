# Paas@Home
This project's aim is to build a service environment using a RaspberryPi and Docker. The documentation includes detailed information about how to set it up and how to run services on it. It should be something like a Paas@Home in the end. I do this, because I think it could grow into a nice collection of tutorial projects on Docker, Golang, Machine Learning, Microservices and so on. Please do not use any of these services in production, just use these as tutorials and howtos on various subjects.
This whole thing is still work in progress, so make sure to come back from time to time to see what's new.

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

### Install Git to get latest service versions
```
sudo apt-get install git
```

## Pull service to server
First create directory to pull the services to and the clone the repository.
```
cd /media/external/
mkdir src
git clone https://github.com/salendron/home-service-host.git
```

## Build and run a service
To make our service available in Docker we need to build them. We also need to comment out some lines from the docker file, since they are only needed by VSCode and not supported by docker on the RaspberryPi. So open the Docker file and comment out or remove these lines.
```
# Install Libs needed for vscode
# RUN go get golang.org/x/tools/gopls
# RUN go get github.com/go-delve/delve/cmd/dlv
```
Navigate to the service directory (the one with the Dockerfile inside) and run **docker build --tag SERVICENAME:VERSION DIRECTORY**.
You have to repeat this process for every service update.
### Build Example:
```
docker build --tag in-memory-db:1.0 .
```
Now we can run our service on docker. We set name, the name of this container instance, as well as the port we want it to run on and also a restart policy. We use "unless-stopped", which restarts the container, even on failures or docker deamon restarts, unless we manually stop it. See more restart options [here](https://docs.docker.com/config/containers/start-containers-automatically/).
```
docker run -d -p 7000:7000 --name in-memory-db -e PORT='7000' -v /var/run/docker.sock:/var/run/docker.sock --restart unless-stopped in-memory-db:1.0
```
You can now use **docker ps** to verify that the service is running.