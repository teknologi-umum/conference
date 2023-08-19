FROM golang:1.21-bookworm AS build

WORKDIR /app

COPY . .

RUN go build -o conf-backend -ldflags="-s -w" .

FROM debian:bookworm-slim AS runtime

WORKDIR /app 

ARG PORT=8080

RUN mkdir -p /app/csv

COPY --from=build /app/ .

EXPOSE ${PORT}

CMD ["/app/conf-backend", "server"]