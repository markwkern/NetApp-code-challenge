#!/bin/bash
theOS=$(uname)
if [[ "$theOS" != "Linux" ]]; then
    echo "This script must be run on a Linux host!"
    exit 1
fi

theDocker=$(which docker)
if [[ -f $theDocker ]]; then
    echo "Found docker at $theDocker"
else
    echo "docker not found, aborting!"
    exit 1
fi
docker run --rm --name imagesrv netappt1:latest &

# is 4 seconds long enough?
sleep 4

theContainer=$(docker ps | grep imagesrv | awk '{ print $1 }')

if [[ -z $theContainer ]]; then
    sleep 5
    theContainer=$(docker ps | grep imagesrv | awk '{ print $1 }')
    if [[ -z $theContainer ]]; then
        echo "The container netappt1:latest did not start."
        echo "You may need to edit this script and increase the sleep time."
        echo "Exiting"
        exit 1
    fi
fi
echo "The container id is: $theContainer"

theAddress=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $theContainer)
if [[ -z $theAddress ]]; then
    echo "The container IP address could not be determined.  Exiting"
    exit 1
fi

echo "Attempting to execute your web browser with $theAddress:9080"

theFire=$(which firefox)
if [[ -f $theFire ]]; then
    $theFire $theAddress:9080
else 
    theXdg=$(which xdg-open)
    if [[ -f $theXdg ]]; then
        theStatus=$($theXdg $theAddress:9080)
        if [[ $theStatus ]]; then
            echo "xdg-open failed (again)"
            echo "Try plugging $theAddress:9080 into your browser URL bar."
        fi
    else
        echo "Could not find web browser"
        echo "Try plugging $theAddress:9080 into your browser URL bar."
    fi
fi

echo "When done, issue the command: docker stop imagesrv"
