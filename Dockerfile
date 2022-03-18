FROM node:gallium-bullseye-slim AS base

WORKDIR /instant

# install curl
RUN apt-get update; apt-get install -y curl

# install docker engine
RUN curl -sSL https://get.docker.com/ | sh

# install docker-compose binary
RUN curl -L "https://github.com/docker/compose/releases/download/1.25.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose

# remove orphan container warning
ENV COMPOSE_IGNORE_ORPHANS=1

# install node deps
ADD package.json .
ADD yarn.lock .
RUN yarn --prod

# add entrypoint script
ADD instant.ts .

ENTRYPOINT [ "yarn", "instant" ]

FROM base as instant-build
# Add default instant OpenHIE packages
ADD . .

RUN apt-get install -y unzip

# install aws cli - for credential fetching
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install

# install kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl