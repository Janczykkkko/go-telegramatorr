FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go get github.com/mattn/go-sqlite3
ENV CGO_ENABLED=1
RUN apk add --no-cache build-base sqlite-dev

RUN go build .

FROM alpine:latest

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /app/telegramatorr .

RUN apk add --no-cache tzdata

EXPOSE 8080

CMD ["/telegramatorr"]