FROM golang:1.14-alpine as builder

WORKDIR /build

ARG VERSION
ARG CI_JOB_TOKEN
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://proxy.golang.org,direct" 

RUN apk add --no-cache git


COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -a -ldflags "-w -X 'main.Version=$VERSION'" -tags 'netgo osusergo' -o /go/bin/ethservice main.go
RUN ldd /go/bin/ethservice 2>&1 | grep -q 'Not a valid dynamic program'

LABEL description=ethservice
LABEL builder=true
LABEL maintainer='Rodrigo <rodrigo@blockwallet.io>'

FROM alpine
COPY --from=builder go/bin/ethservice /usr/local/bin

WORKDIR usr/local/bin
ENTRYPOINT [ "ethservice", "serve" ]
EXPOSE 8080 8443 9008
