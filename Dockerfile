# Build backend
FROM golang:1.24.0-alpine3.20 as build_backend

RUN apk add upx alpine-sdk

WORKDIR /backend

COPY [^ui] .

RUN go mod download

RUN CGO_ENABLED=1 go build -ldflags "-s -w" -v -tags musl cmd/main/main.go

RUN upx --best --lzma main


# Build frontend
FROM node:22.8.0-alpine3.20 as build_frontend

WORKDIR /frontend

COPY uit/package.json package.json

COPY ui/pnpm-lock.yaml pnpn-lock.yaml

RUN npm install -g pnpm@9.15.5
RUN pnpm install

COPY ui/[^node_modules] .

RUN pnpm run build


# End container
FROM alpine:3.20

WORKDIR /

COPY --from=build_backend /backend/main .
COPY --from=build_frontend /frontend/dist public

ENV ENV=PRODUCTION

EXPOSE 4000

ENTRYPOINT ["./main"]
