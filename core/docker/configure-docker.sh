#!/bin/bash

echo "Configuring docker"

sudo groupadd docker
sudo usermod -aG docker $USER

if [ "$?" == "1" ]; then
    echo "Could not add user to docker group"
else
    echo 'User added to docker group successfully'
fi

sudo chown "$USER":"$USER" /home/"$USER"/.docker -R
sudo chmod g+rwx "$HOME/.docker" -R
newgrp docker

echo "Finished configuring."
