FROM golang:latest

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go mod tidy

RUN go build -o main ./cmd

EXPOSE 8080

CMD ["./main"]