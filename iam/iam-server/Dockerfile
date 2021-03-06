FROM golang:1.16 as builder

WORKDIR /workspace

# Get the dependencies so it can be cached into a layer
COPY go.mod go.sum ./
RUN go mod download

# Now copy all the source...
COPY . .

ARG revisionID=unknown
ARG buildTimestamp=unknown

# ...and build it.
RUN CGO_ENABLED=0 go build -o ./iam/bin/iam-server \
    -ldflags="-s -w -X main.revisionID=${revisionID} -X main.buildTimestamp=${buildTimestamp}" \
    ./iam/iam-server

# Build the runtime image
FROM alpine:3.13
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Include resources required by IAM server
COPY ./iam/iam-server/resources ./resources

COPY --from=builder /workspace/iam/bin/iam-server ./server

# Web UI
EXPOSE 8080
# REST
EXPOSE 9080
# gRPC
EXPOSE 50051

ENTRYPOINT ["./server"]
