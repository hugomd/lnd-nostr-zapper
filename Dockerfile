FROM golang:1.20-alpine AS build
ENV GO111MODULE=on
WORKDIR /go/src/github.com/hugomd/lnd-nostr-zapper/
COPY . /go/src/github.com/hugomd/lnd-nostr-zapper/
RUN cd /go/src/github.com/hugomd/lnd-nostr-zapper && \
    go get && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.18.0
COPY --from=build /go/src/github.com/hugomd/lnd-nostr-zapper/main /
COPY entrypoint.sh /entrypoint.sh
RUN apk add ca-certificates tor
ENTRYPOINT ["/entrypoint.sh"]
