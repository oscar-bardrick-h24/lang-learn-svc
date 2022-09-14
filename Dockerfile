FROM golang:1.18-alpine

EXPOSE 8080

RUN apk add make

COPY . /lang-learn-svc
WORKDIR /lang-learn-svc

RUN go build cmd/main.go

CMD ["./main"]