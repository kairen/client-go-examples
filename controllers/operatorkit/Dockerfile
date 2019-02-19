FROM golang:1.9-alpine AS build-env
MAINTAINER Kyle Bai(kyle.b@inwinstack.com)

ENV GOPATH "/go"
RUN apk add --no-cache make git g++
ADD . /go/src/github.com/kairen/simple-operator
RUN cd /go/src/github.com/kairen/simple-operator && \
  make deps && \
  IN_DOCKER=1 make

# Run stage
FROM alpine

COPY --from=build-env /go/src/github.com/kairen/simple-operator/out/operator /bin/operator
RUN apk add --no-cache ca-certificates && \
  rm -f /var/cache/apk/*

ENTRYPOINT ["/bin/operator"]
