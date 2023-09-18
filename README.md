# GOBANK
Simple example of creating an API with GO while using minimal amount of libaries

# HOW TO USE
## SET THE JWT_SECRET
Before starting must set the environment variable JWT_SECRET by typing `export JWT_SECRET="your_secret_goes_in_here"`
## MAKE POSTGRES AND PGADMIN4 WORK WITH USING DOCKER
Make sure docker is installed then **cd** into the folder where the **docker-compose.yml** is located
Run `docker compose up -d` to pull and start the containers

## CONFIGURING PGADMIN4 AND POSTGRES
Open the browser and type in **localhost:5050**
Login in with the credentials provided to the container inside the **docker-compose.yml** for pgadmin
Than add a server with the credentials provided inside the **docker-compose.yml** for db
For the filed **IP-Adress** type in the name of the postgres container **postgres_container**
Make sure to disable SSL for the moment as its not configured at the moment