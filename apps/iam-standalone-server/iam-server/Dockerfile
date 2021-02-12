FROM golang:1.15 as builder

WORKDIR /workspace

# Get the dependencies so it can be cached into a layer
COPY iam/go.mod iam/go.sum ./iam/

# Override with local modules
COPY foundation foundation
RUN echo "replace github.com/kadisoka/foundation => ../foundation\n" >> iam/go.mod

WORKDIR /workspace/iam

RUN go mod download

# Now copy all the source...
COPY ./iam/ .

ARG revisionID=unknown
ARG buildTimestamp=unknown

RUN env

# ...and build it.
RUN CGO_ENABLED=0 go build -o ./bin/iam-server \
    -ldflags="-s -w -X main.revisionID=${revisionID} -X main.buildTimestamp=${buildTimestamp}" \
    ./iam-server

# Build the runtime image
FROM alpine:3.13
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Include resources required by IAM server
COPY ./iam/iam-server/resources ./resources

COPY --from=builder /workspace/iam/bin/iam-server ./service

# Web UI
EXPOSE 8080
# REST
EXPOSE 9080
# gRPC
EXPOSE 50051

ENTRYPOINT ["./service"]
