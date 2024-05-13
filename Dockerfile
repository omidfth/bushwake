# base go image
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app/cmd/api

RUN export GOPROXY=https://goproxy.io,direct && CGO_ENABLED=0 go build -o bushwake .

RUN chmod +x /app/cmd/api/bushwake

# build a tiny docker image
FROM alpine:latest

RUN apk add tzdata

RUN mkdir /app

COPY --from=builder /app/cmd/api/bushwake /app

COPY config.yml /app

WORKDIR /app

CMD [ "/app/bushwake","-e","/app/config.yml"]