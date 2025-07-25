FROM golang:1.24-alpine AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o auth_test ./cmd/authserv/main.go

CMD ["./auth_test"]