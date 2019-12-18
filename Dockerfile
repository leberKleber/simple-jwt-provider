# Go build
FROM golang as build

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOPATH=/

WORKDIR /src/simple-jwt-provider/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags -s -a -o simple-jwt-provider -installsuffix cgo ./cmd/provider/

# Service definition
FROM alpine

RUN apk add --update libcap tzdata ca-certificates && rm -rf /var/cache/apk/*

COPY --from=build /src/simple-jwt-provider/simple-jwt-provider simple-jwt-provider

RUN setcap CAP_NET_BIND_SERVICE=+eip ./simple-jwt-provider
RUN update-ca-certificates

RUN addgroup -g 1000 -S runnergroup && adduser -u 1001 -S apprunner -G runnergroup
USER apprunner

ENTRYPOINT ["/simple-jwt-provider"]
