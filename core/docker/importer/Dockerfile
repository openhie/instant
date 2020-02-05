FROM node:erbium-alpine

WORKDIR /importer

RUN npm install -g wait-on

# copy all files in this directory into the docker image
COPY . .

# this file isnt needed as we are already dockerized
RUN rm Dockerfile
