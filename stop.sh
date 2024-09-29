#!/bin/bash
docker-compose -f docker-compose-orderer.yaml -f docker-compose-peer-couchdb.yaml down

docker volume prune
