# syntax=docker/dockerfile:1@sha256:4edf897a3ffa55b89f906fc8cc78afdb3f1834cc9c7083565e611a8a7d5fe99e
FROM node:24-alpine3.24@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53

ENV WEBUI_DIR=/src/webui
RUN mkdir -p $WEBUI_DIR

COPY package.json yarn.lock .yarnrc.yml $WEBUI_DIR/

ENV VITE_APP_BASE_URL=""
ENV VITE_APP_BASE_API_URL="/api"

WORKDIR $WEBUI_DIR

RUN npm i -g corepack
RUN yarn workspaces focus --all --production

COPY . $WEBUI_DIR/

EXPOSE 8080
