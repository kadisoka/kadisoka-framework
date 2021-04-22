FROM golang:1.16

RUN go get -u golang.org/x/lint/golint

ENTRYPOINT ["golint"]
