FROM golang:1.20-alpine AS build
ENV GO111MODULE=on
WORKDIR /go/src/github.com/hugomd/lnd-nostr-zapper/
COPY . /go/src/github.com/hugomd/lnd-nostr-zapper/
RUN cd /go/src/github.com/hugomd/lnd-nostr-zapper && \
    go get && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.18.0
WORKDIR /golang
RUN adduser -D golang -h /golang && apk add ca-certificates tor
USER golang
COPY --from=build --chown=golang:0 /go/src/github.com/hugomd/lnd-nostr-zapper/main /golang
COPY --chown=golang:0 entrypoint.sh /golang/entrypoint.sh
ENTRYPOINT ["/golang/entrypoint.sh"]
