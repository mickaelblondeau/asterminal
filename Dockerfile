FROM golang:1.26-alpine AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

ARG VERSION=dev

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/asterminal cmd/main.go

FROM alpine:3.23.4 AS runtime

COPY --from=build /usr/local/bin/asterminal /usr/local/bin/asterminal

ENTRYPOINT ["/usr/local/bin/asterminal"]
