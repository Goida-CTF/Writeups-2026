##### image tags #####
ARG GO_BUILD_TAG=1.25.6-alpine3.22
ARG ALPINE_TAG=3.22

###### build #####
FROM golang:${GO_BUILD_TAG} AS build

WORKDIR /build/

COPY go.mod ./

COPY cmd/ ./cmd/

RUN go build -v -o server ./cmd/server/

##### runtime #####
FROM build AS artifact-source
FROM alpine:${ALPINE_TAG}

WORKDIR /app/

COPY --from=artifact-source /build/server ./

RUN chmod +x /app/server

EXPOSE 5010

ENTRYPOINT [ "/app/server" ]
