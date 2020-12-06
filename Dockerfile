FROM centos:latest

COPY ./Dockerfile /dockerfile/
COPY ./terway-probe /usr/local/bin/terway-probe
CMD bash