
Within this directory you will find a example docker application which serves up a directory of images
supplied in the directory tree_for_candidate in a web browser.  Some nominal information
about the images files is provided along with a link to view the images.

Prerequisites: 
    1) docker 
    2) A suitable web browser which understands xdg-open, e.g. firefox or chrome.
    3) If desired to build the application, golang version 1.16 or higher

The golang application is in the subdirectory named servit.  To build, invoke 
the bash script build-hack.sh therein and the executable will land in
build/imagesrv.  This executable has already been copied to the subdirectory
docker-setup.

To build the container, invoke the bash script build-container.sh.

To start the resulting container image, invoke the bash script invokeit.sh.

If no suitable web browser is found, an IP address and port are emitted from invokeit.sh.
Plug this address:port into your web browser's URL bar.

Note the assumption is made that the web browser is running on the same Linux machine (or
VM) where the container is running.


Mark W. Kernodle
mark.kernodle@gmail.com

