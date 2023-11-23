#!/bin/bash


# Stop the running container
echo -e "Stopping containers...\n"

sudo docker stop Forum


# Remove the stopped container
echo -e "Removing stopped container...\n"
sudo docker rm Forum

# Remove the images
echo -e "Removing images...\n"

sudo docker rmi forum_image
