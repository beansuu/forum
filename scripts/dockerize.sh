# !/bin/bash

# Build the images
echo -e "Building images...\n"

docker build -f Dockerfile -t forum_image .

# To run docker container

echo -e "Running docker container...\n"

docker run -dp 8000:8000 --name Forum forum_image