# CONTAINER FOR BUILDING BINARY
FROM golang:1.22.4 AS build

WORKDIR $GOPATH/src/github.com/0xPolygon/cdk

# INSTALL DEPENDENCIES
COPY go.mod go.sum ./
RUN go mod download

# BUILD BINARY
COPY . .
RUN make build-go

# CONTAINER FOR RUNNING BINARY
FROM alpine:3.18.4

COPY --from=build /go/src/github.com/0xPolygon/cdk/target/cdk-node /usr/local/bin/

CMD ["/bin/sh", "-c", "cdk"]
