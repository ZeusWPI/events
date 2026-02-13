# Stage 1: Build backend
FROM golang:1.26.0-alpine3.22 AS backend-builder
WORKDIR /app

RUN apk add --no-cache gcc musl-dev libwebp-dev

COPY go.mod go.sum ./ 
RUN go mod download 

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o migrate migrate.go


# Stage 3: Build frontend
FROM node:22.16.0-alpine3.22 AS frontend-builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /frontend/ui

COPY ./ui/package.json ./ui/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./ui/ .

ARG BUILD_TIME

ENV VITE_BUILD_TIME=$BUILD_TIME
ENV CI=true

RUN pnpm run build

# Stage 3: Combine
FROM alpine:3.22 AS prod
WORKDIR /app

RUN apk add --no-cache libwebp-dev ca-certificates tzdata 

COPY --from=backend-builder /app/server .
COPY --from=backend-builder /app/migrate .
COPY --from=backend-builder /app/docs ./docs
COPY --from=frontend-builder /frontend/ui/dist ./public

ENV APP_ENV=PRODUCTION

EXPOSE 4000

ENTRYPOINT ["sh", "-c", "./migrate && ./server"]
