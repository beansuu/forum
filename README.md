# Forum

## Description
This project aims to create a web forum with the following key features:

### Communication Between Users
Users can create posts and comments.

### Likes and Dislikes
Only registered users can like or dislike posts and comments.
The number of likes and dislikes is visible to all users.

### Filter Mechanism
Users can filter displayed posts by categories, created posts, and liked posts.
Filtering by categories is akin to subforums.

### Authentication
Users can register by providing their email, username, and password.

### SQLite Database
Data, including users, posts, and comments, is stored using the SQLite database.

### Docker Integration
The project is containerized using Docker for easy deployment.
Basic Docker knowledge is recommended; refer to the provided Docker basics resource.

## Usage
### Using Docker
1.  Clone the repository

2. Build Docker image: <code>docker image build -f Dockerfile -t <name_of_the_image> .</code> Example: <code>docker image build -f Dockerfile -t ascii-art-web1 .</code>

3. Start the container: docker container <code>run -p <port_you_what_to_run> --detach --name <name_of_the_container> <name_of_the_image></code> Example: <code>docker container run -p 8070:8070 --detach --name dockerized-ascii-art ascii-art-web1</code>

4. Run the container to start server: <code>docker exec -it <container_name> /app/bin Example: docker exec -it dockerized-ascii-art /app/bin</code>

### For local testing:

**Make sure you have Docker running**

```
> bash scripts/dockerize.sh 
```

After completing tests, run the following to clean up images
```
> bash scripts/cleanup.sh 
```

## Authors
### Mauno Talli 
<a href="https://01.kood.tech/git/mtalli">@mtalli</a>

### Andreas Selge 
<a href="https://01.kood.tech/git/aselge">@aselge</a>