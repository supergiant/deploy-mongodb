FROM ubuntu:14.04
MAINTAINER Qbox Inc.

RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
RUN echo "deb http://repo.mongodb.org/apt/ubuntu trusty/mongodb-org/3.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-3.2.list
RUN apt-get update
RUN apt-get install -y mongodb-org-shell

ADD deploy-mongodb deploy-mongodb
ENTRYPOINT ["/deploy-mongodb"]
