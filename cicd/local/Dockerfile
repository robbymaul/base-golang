FROM golang:1.23.0-alpine AS build

# Install tzdata
RUN apk add --no-cache make

# Define arguments
ARG ARG_SERVICE_SLUG
ARG ARG_APP_VERSION

# Set working directory
WORKDIR /usr/src/app

# Copy files
COPY .git ./.git
COPY api ./api
COPY config ./config
COPY pkg ./pkg
COPY init ./init
COPY go.mod go.sum Makefile ./
# -- Copy service specific files
COPY cmd/${ARG_SERVICE_SLUG} ./cmd/${ARG_SERVICE_SLUG}
COPY app/ ./app/

ENV PROJECT_SLUG=${ARG_SERVICE_SLUG} \
    APP_VERSION=${ARG_APP_VERSION}

RUN make release-bin

FROM alpine:3.20

# Install tzdata
RUN apk add --no-cache tzdata

# Create a non-root user
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

# Define arguments
ARG ARG_PORT=10000
ARG ARG_SERVICE_SLUG

COPY --from=build /usr/src/app/bin/release/${ARG_SERVICE_SLUG} /usr/local/bin/svc
COPY migrations /etc/svc/migrations

EXPOSE ${ARG_PORT}

ENV SERVICE_SLUG=${ARG_SERVICE_SLUG}

USER nonroot

ENTRYPOINT ["svc", "-dir", "/etc/svc"]
