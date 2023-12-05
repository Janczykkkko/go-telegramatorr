FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build .

FROM scratch 
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /app/telegramatorr .

CMD ["/telegramatorr"]