#!/bin/bash

export PATH=$PATH:/usr/bin
cd /tmp/

docker compose -f docker_compose.yml down 
docker compose -f docker_compose.yml up -d 
docker compose -f docker_compose.yml pull