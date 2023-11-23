#!/bin/bash


# Stop the running container
echo -e "Stopping containers...\n"

docker stop Forum


# Remove the stopped container
echo -e "Removing stopped container...\n"
docker rm Forum

# Remove the images
echo -e "Removing images...\n"

docker rmi forum_image
