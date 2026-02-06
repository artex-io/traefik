# -- WEBUI ---------------------------------------------------------------------

FROM --platform=$BUILDPLATFORM node:24.7-alpine3.22 AS webui

RUN npm upgrade --global yarn

WORKDIR /src/webui/

COPY ./webui/package.json ./webui/yarn.lock ./webui/.yarnrc.yml ./

RUN yarn install

COPY ./webui/ ./

RUN yarn build

# -- GO BUILD ------------------------------------------------------------------

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS gobuild

WORKDIR /go/src/github.com/traefik/traefik

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN apk --update upgrade \
    && apk --no-cache --no-progress add make git mercurial bash gcc musl-dev curl tar ca-certificates tzdata libcap \
    && update-ca-certificates

COPY . .

RUN rm -rf static/

COPY --from=webui /src/webui/static/ ./webui/static/

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

SHELL ["bash", "-c"]

RUN if [ "${TARGETARCH}" = "amd64" ]; then \
        VERSION="$(git describe --tags --always)" GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOAMD64=${TARGETVARIANT} make binary; \
    elif [ "${TARGETARCH}" = "arm" ]; then \
        VERSION="$(git describe --tags --always)" GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT} make binary; \
    elif [ "${TARGETARCH}" = "arm64" ]; then \
        VERSION="$(git describe --tags --always)" GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM64=${TARGETVARIANT} make binary; \
    else \
        echo "Unsupported architecture: ${TARGETARCH}"; exit 1; \
    fi

RUN setcap cap_net_bind_service=+ep "dist/$TARGETPLATFORM/traefik"

# -- scratch -------------------------------------------------------------------

FROM scratch

ARG TARGETPLATFORM

COPY --from=gobuild /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuild /etc/passwd /etc/passwd
COPY --from=gobuild /etc/group /etc/group
COPY --from=gobuild /etc/services /etc/services

COPY --from=gobuild /go/src/github.com/traefik/traefik/dist/${TARGETPLATFORM}/traefik /

EXPOSE 80
VOLUME ["/tmp"]

ENTRYPOINT ["/traefik"]
