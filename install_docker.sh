#!/bin/bash
sudo apt update
sudo apt install --yes apt-transport-https ca-certificates curl gnupg2 software-properties-common git
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
sudo apt update
sudo apt install --yes docker-ce
sudo usermod -aG docker $USER
newgrp docker
echo "Docker installation and setup complete."
