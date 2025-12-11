FROM golang:1.24-alpine

RUN apk add --no-cache git bash

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -buildvcs=false -o hh ./cmd/hh
RUN chmod +x hh

EXPOSE 8080

CMD ["CompileDaemon", "--build=go build -buildvcs=false -o hh ./cmd/hh", "--command=./hh", "--directory=./", "--include=cmd/hh|internal", "--exclude=vendor|\\.git"]
