FROM node:gallium-bullseye-slim

WORKDIR /instant

# install node deps
ADD package.json .
ADD yarn.lock .
RUN yarn

ADD test.sh .

ENTRYPOINT [ "yarn", "test:container" ]
