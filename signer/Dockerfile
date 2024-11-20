# CONTAINER FOR BUILDING BINARY
FROM golang:1.22 AS build

# INSTALL DEPENDENCIES
RUN go install github.com/gobuffalo/packr/v2/packr2@v2.8.3
COPY go.mod go.sum /src/
RUN cd /src && go mod download

# BUILD BINARY
COPY . /src
RUN cd /src/db && packr2
RUN cd /src/signer && make build

# CONTAINER FOR RUNNING BINARY
FROM alpine:3.18.0
COPY --from=build /src/signer/dist/xlayer-signer /app/xlayer-signer

EXPOSE 7001
CMD ["/bin/sh", "-c", "/app/xlayer-signer"]
