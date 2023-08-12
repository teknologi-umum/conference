FROM golang:1.21-bookworm AS build

WORKDIR /app

COPY . .

RUN go build -o app -ldflags="-s -w" .

FROM debian:bookworm-slim AS runtime

WORKDIR /app 

ARG PORT=8080

COPY --from=build /app/ .

EXPOSE ${PORT}

CMD ["./app"]