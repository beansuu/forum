# !/bin/bash

# Build the images
echo -e "Building images...\n"

sudo docker build -f Dockerfile -t forum_image .

# To run docker container

echo -e "Running docker container...\n"

sudo docker run -dp 8000:8000 --name Forum forum_image

echo -e "Open browser and go to http://localhost:8000...\n"