##### image tags #####
ARG GO_BUILD_TAG=1.23-alpine
ARG TARGET_TAG=3.21

##### build #####
FROM golang:${GO_BUILD_TAG} AS build

WORKDIR /build/

COPY go.mod go.sum ./
RUN go mod download

RUN apk add --no-cache build-base

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -v -o initial ./cmd/initial \
    && go build -v -o server ./cmd/server

##### runtime #####
FROM build AS artifact-source
FROM alpine:${TARGET_TAG}

WORKDIR /app/

COPY --from=artifact-source /build/initial ./
COPY --from=artifact-source /build/server ./
COPY ./docker/backend/docker-entrypoint.sh ./

RUN chmod +x /app/initial /app/server \
    /app/docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT [ "/app/docker-entrypoint.sh" ]
