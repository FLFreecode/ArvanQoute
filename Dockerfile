FROM golang:1.21 AS builder

ADD . /tmp/app
WORKDIR /tmp/app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o arvan-qoute main.go

FROM alpine 

#RUN apk --no-cache

WORKDIR /app

COPY --from=builder /tmp/app/qoute-service .
COPY --from=builder /tmp/app/config.yml .

EXPOSE 8690

RUN chmod +x ./qoute-service

CMD ["./arvan-qoute"]
