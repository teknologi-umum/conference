FROM golang:1.21-bookworm AS build

WORKDIR /app

COPY . .

RUN go build -o conf-backend -ldflags="-s -w -X main.version=$(git rev-parse HEAD)" .

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /app/csv && \
    mkdir -p /data

ARG PORT=8080

COPY --from=build /app/ .

HEALTHCHECK --interval=60s --timeout=40s \
    CMD /app/conf-backend --port ${PORT}

EXPOSE ${PORT}

CMD ["/app/conf-backend", "server"]