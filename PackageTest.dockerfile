FROM node:fermium-buster

WORKDIR /instant

# install node deps
ADD package.json .
ADD yarn.lock .
RUN yarn

ADD test.sh .

ENTRYPOINT [ "yarn", "test:container" ]
