# Multi-stage build
FROM golang:1.22-alpine AS build
WORKDIR /src


COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download


COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/ems-proxy ./cmd/ems-proxy


# Minimal runtime image
FROM scratch
COPY --from=build /out/ems-proxy /ems-proxy
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/ems-proxy"]