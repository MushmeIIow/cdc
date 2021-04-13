FROM golang:1.15-alpine3.12 AS build

RUN apk add --no-cache git

WORKDIR /go/src/github.com/mushmellow/cdc-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o app ./cmd/main.go ./cmd/init.go

FROM alpine:3.12

WORKDIR /root/
COPY --from=build /go/src/github.com/mushmellow/cdc-service/app .
CMD ["./app"]