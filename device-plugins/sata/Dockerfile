FROM golang:1.9-alpine AS build-env
MAINTAINER Kyle Bai(k2r2bai@gmail.com)

ENV GOPATH "/go"
RUN apk add --no-cache make git g++ && \
  go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/kairen/simple-device-plugin
RUN cd /go/src/github.com/kairen/simple-device-plugin && \
  dep ensure && IN_DOCKER=1 make

# Run stage
FROM alpine:3.7

COPY --from=build-env /go/src/github.com/kairen/simple-device-plugin/out/sdp /bin/sdp
RUN apk add --no-cache util-linux bash && \
  rm -f /var/cache/apk/*

ENTRYPOINT ["/bin/sdp"]
