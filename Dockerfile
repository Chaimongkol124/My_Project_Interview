FROM golang:1.22.0-alpine3.19

WORKDIR /app

COPY . .

RUN go build -o api

EXPOSE 8080
CMD ["./api"]