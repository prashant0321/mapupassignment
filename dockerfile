# Dockerfile
FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod init users/91966/mapupassignment
RUN go build -o app.go 

CMD ["./main"]
