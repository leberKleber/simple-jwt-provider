FROM golang:latest

RUN mkdir -p /go/mods/tests
WORKDIR /go/mods/tests

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . /go/mods/tests

ENV RUN_TESTS=""
CMD go test -v -count=1 -tags=component ${RUN_TESTS} ./cmd/provider