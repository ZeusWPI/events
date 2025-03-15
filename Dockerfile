# Build backend
FROM golang:1.24.0-alpine3.20 as build_backend

RUN apk add upx alpine-sdk

WORKDIR /backend

COPY go.mod go.sum ./ 
RUN go mod download 

COPY . .

# Build server executable
RUN CGO_ENABLED=1 go build -ldflags "-s -w" -v -tags musl -o main ./cmd/events/main.go
RUN upx --best --lzma main

# Build migration executable
RUN CGO_ENABLED=1 go build -ldflags "-s -w" -v -tags musl -o migrate migrate.go
RUN upx --best --lzma migrate


# Build frontend
FROM node:22.8.0-alpine3.20 as build_frontend

WORKDIR /frontend

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN npm install -g pnpm@9.15.5 && pnpm install 

COPY ui/ .

ARG BUILD_TIME
ENV VITE_BUILD_TIME=$BUILD_TIME

RUN pnpm run build


# End container
FROM alpine:3.20

WORKDIR /

COPY --from=build_backend /backend/main .
COPY --from=build_backend /backend/migrate .
COPY --from=build_frontend /frontend/dist ./public

RUN chmod +x ./main ./migrate

ENV ENV=PRODUCTION

EXPOSE 4000

ENTRYPOINT ["sh", "-c", "./migrate && exec ./main"]
