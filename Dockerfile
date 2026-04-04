FROM node:22-bookworm-slim AS frontend-build

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build


FROM golang:1.25-bookworm AS backend-build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server


FROM debian:bookworm-slim

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN mkdir -p /app/uniswap-api /app/frontend

COPY --from=backend-build /out/server ./server
COPY registry.yaml ./registry.yaml
COPY ebooks ./ebooks
COPY uniswap-api/config.yaml ./uniswap-api/config.yaml
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

EXPOSE 8080

CMD ["./server"]
