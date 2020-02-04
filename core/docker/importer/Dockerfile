FROM node:erbium-alpine

WORKDIR /importer

COPY package.json package-lock.json ./

RUN npm install

# copy all files in this directory into the docker image
COPY . .

# this file isnt needed as we are already dockerized
RUN rm Dockerfile
