FROM golang:1.23.4 AS build

WORKDIR /app-build

COPY go.mod go.sum ./

RUN go get ./...

COPY . .

RUN go build -o ./build/app ./cmd/server/

FROM ubuntu:25.04 AS local-app

WORKDIR /app

COPY --from=build /app-build/build/app .

EXPOSE 1323
EXPOSE 50051

ENV DB_HOST="192.168.2.44"

ENTRYPOINT ["./app"]

FROM ubuntu:25.04 AS app

WORKDIR /app

COPY --from=build /app-build/build/app .

EXPOSE 1323
EXPOSE 50051

ENTRYPOINT ["./app"]
