FROM golang:1.23.4 AS build

WORKDIR /app-build

COPY go.mod go.sum ./

RUN go get ./...

COPY ./cmd ./internal ./pkg ./proto ./

RUN go build -o ./build/app ./cmd/server

FROM ubuntu:25.04 AS app

WORKDIR /app

COPY --from=build /build/app .

ENTRYPOINT ["./app"]
