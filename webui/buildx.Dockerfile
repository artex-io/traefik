FROM node:24-alpine3.22@sha256:191c9f0080fcbbc6547a85dc0ff7988072214a355aabdc1d2ec55a7dae5eea8a

ENV WEBUI_DIR=/src/webui
RUN mkdir -p $WEBUI_DIR

COPY package.json yarn.lock .yarnrc.yml $WEBUI_DIR/

ENV VITE_APP_BASE_URL=""
ENV VITE_APP_BASE_API_URL="/api"

WORKDIR $WEBUI_DIR

RUN corepack enable
RUN yarn workspaces focus --all --production

COPY . $WEBUI_DIR/

EXPOSE 8080
