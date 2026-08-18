# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/eveping ./cmd/eveping

FROM alpine:3.20
RUN adduser -D -u 10001 eveping
COPY --from=build /out/eveping /usr/local/bin/eveping
USER eveping
ENTRYPOINT ["/usr/local/bin/eveping"]
