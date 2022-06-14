FROM golang:1.18

RUN go get -u golang.org/x/lint/golint

ENTRYPOINT ["golint"]
