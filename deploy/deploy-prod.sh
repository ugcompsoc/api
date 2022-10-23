#!/bin/bash

# Stop the docker container currently running the API
echo -e "\n > Shutting down CompSoc API"
docker-compose -f docker-compose.prod.yml down

# Pull down any new workflow related stuff
echo -e "\n > Pulling down new updates from main"
git checkout main
git pull 

# Start the service again
echo -e "\n > Starting CompSoc API"
docker-compose -f docker-compose.prod.yml up -d
