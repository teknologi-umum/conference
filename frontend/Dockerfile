FROM node:20.5.1-bookworm AS builder

WORKDIR /app

COPY . .

RUN npm i -g pnpm@latest && \
    pnpm install && \
    pnpm run build

FROM node:20.5.1-bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* 

COPY --from=builder /app/.output .

HEALTHCHECK --interval=60s --timeout=40s \
  CMD curl -f http://localhost:3000/ || exit 1

EXPOSE 3000

CMD ["node", "./server/index.mjs"]
