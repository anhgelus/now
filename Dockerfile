FROM node:22 as builder

WORKDIR /app

COPY . .

RUN npm install -g bun

RUN bun install && bun run build

FROM golang:1.23-alpine

WORKDIR /app

COPY --from=builder . .

RUN go mod tidy && go mod build -o app .

ENV NOW_DOMAIN=""
ENV NOW_DATA=""

EXPOSE 80

CMD ./app -domain $NOW_DOMAIN -data $NOW_DATA
