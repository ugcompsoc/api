#!/bin/bash

# Change version of API
if $1 -eq ""
then
    exit 1
fi
export API_VERSION=$1

# Stop the docker container currently running the API
echo -e "\n > Shutting down CompSoc API"
docker-compose -f docker-compose.prod.yml down

# Start the service again
echo -e "\n > Starting CompSoc API"
docker-compose -f docker-compose.prod.yml up -d
