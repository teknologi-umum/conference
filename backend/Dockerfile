FROM golang:1.21-bookworm AS build

WORKDIR /app

COPY . .

RUN go build -o conf-backend -ldflags="-s -w -X main.version=$(git rev-parse HEAD)" .

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ARG PORT=8080

RUN mkdir -p /app/csv

COPY --from=build /app/ .

HEALTHCHECK --interval=60s --timeout=40s \
    CMD curl -f http://localhost:8080/ping || exit 1

EXPOSE ${PORT}

CMD ["/app/conf-backend", "server"]