#!/bin/bash
theOS=$(uname)
if [[ "$theOS" != "Linux" ]]; then
    echo "This script must be run on a Linux host!"
    exit 1
fi

theDocker=$(which docker)
if [[ -f $theDocker ]]; then
    echo "Found docker at $theDocker, proceeding"
else
    echo "docker not found, aborting!"
    exit 1
fi

if [[ -e ./docker-setup ]]; then
    docker build ./docker-setup -t netappt1:latest
else
    echo "This script must be run in the proscribed project directory!"
fi
