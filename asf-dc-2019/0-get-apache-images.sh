#!/bin/bash

# Example for the Docker Hub V2 API
# Returns all images and tags associated with a Docker Hub organization account.
# Requires 'jq': https://stedolan.github.io/jq/

ORG="$1"

echo "Scanning org $ORG"

set -e

# get list of repositories

REPO_LIST_OFFICIAL=($(docker search $ORG | cut -d' ' -f 1))
REPO_LIST_APACHE=($(curl -k https://hub.docker.com/v2/repositories/${ORG}/?page_size=100 | jq -r '.results|.[]|.name'))
echo "-------------------------"
echo $REPO_LIST_OFFICIAL
echo "-------------------------"

for i in ${REPO_LIST_OFFICIAL[@]}
do
  if [[ $i == "$" ]]; then
	echo ""		
  else
  	echo "image: ${i}:latest"
  fi
done
