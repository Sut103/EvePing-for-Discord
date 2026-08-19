# syntax=docker/dockerfile:1
FROM golang:1.24.7-alpine@sha256:fc2cff6625f3c1c92e6c85938ac5bd09034ad0d4bc2dfb08278020b68540dbb5 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/eveping ./cmd/eveping

FROM alpine:3.20@sha256:d9e853e87e55526f6b2917df91a2115c36dd7c696a35be12163d44e6e2a4b6bc
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 eveping
COPY --from=build /out/eveping /usr/local/bin/eveping
USER eveping
ENTRYPOINT ["/usr/local/bin/eveping"]
