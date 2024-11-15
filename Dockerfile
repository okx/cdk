# CONTAINER FOR BUILDING BINARY
FROM --platform=${BUILDPLATFORM} golang:1.22.4 AS build

WORKDIR $GOPATH/src/github.com/0xPolygon/cdk

# INSTALL DEPENDENCIES
COPY go.mod go.sum ./
RUN go mod download

# BUILD BINARY
COPY . .
RUN make build-go

# CONTAINER FOR RUNNING BINARY
FROM --platform=${BUILDPLATFORM} debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates postgresql-client libssl-dev && rm -rf /var/lib/apt/lists/*
COPY --from=build /go/src/github.com/0xPolygon/cdk/target/cdk-node /usr/local/bin/

CMD ["/bin/sh", "-c", "cdk"]
