FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o network_test

EXPOSE 8080

CMD ["./network_test"]
