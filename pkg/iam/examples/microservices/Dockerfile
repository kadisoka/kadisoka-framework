FROM golang:1.18 as builder

WORKDIR /workspace

# Get the dependencies so it can be cached into a layer
COPY go.mod go.sum ./
RUN go mod download

# Now copy all the source...
COPY . .

ARG revisionID=unknown
ARG buildTimestamp=unknown

# ...and build it.
RUN CGO_ENABLED=0 go build -o ./iam/bin/microservices \
    -ldflags="-s -w -X main.revisionID=${revisionID} -X main.buildTimestamp=${buildTimestamp}" \
    ./iam/examples/microservices

# Build the runtime image
FROM alpine:3.16
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Include resources required by IAM server
COPY ./pkg/iam/iam-server-app/resources ./resources

COPY --from=builder /workspace/iam/bin/microservices ./server

# HTTP
EXPOSE 8080

ENTRYPOINT ["./server"]
