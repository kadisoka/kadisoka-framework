FROM golang:1.16

WORKDIR /workspace

# Get the dependencies so it can be cached into a layer
COPY go.mod go.sum ./
RUN go mod download

ENTRYPOINT [ "go" ]
