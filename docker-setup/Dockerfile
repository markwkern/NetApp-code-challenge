FROM fedora:latest

COPY tree_for_candidate /var/tmp/Netapp-images
COPY imagesrv /usr/local/bin
RUN dnf -y install file
RUN dnf -y install tree
EXPOSE 9080
CMD ["/usr/local/bin/imagesrv", "/var/tmp/Netapp-images"]
