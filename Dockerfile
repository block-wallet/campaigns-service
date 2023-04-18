FROM golang:1.20-alpine as builder

## install gcc tool
RUN apk add build-base

WORKDIR /build

ARG VERSION
ARG CI_JOB_TOKEN
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY="https://proxy.golang.org,direct" 

RUN apk add --no-cache git


COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -a -ldflags "-w -X 'main.Version=$VERSION'" -tags 'netgo osusergo' -o /go/bin/campaignsservice main.go
# RUN ldd /go/bin/campaignsservice 2>&1 | grep -q 'Not a valid dynamic program'

LABEL description=campaignsservice
LABEL builder=true
LABEL maintainer='Facundo <facundo@blockwallet.io>'

FROM alpine
COPY --from=builder go/bin/campaignsservice /usr/local/bin

WORKDIR /usr/local/bin
ENTRYPOINT [ "campaignsservice", "serve" ]
EXPOSE 8080 8443 9008
