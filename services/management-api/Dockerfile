FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

#
RUN go build -o management-api ./cmd/server

EXPOSE 81

CMD ["./management-api"]
