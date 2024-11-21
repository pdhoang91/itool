# Giai đoạn build
FROM golang:1.21 AS builder

WORKDIR /app

# Cài đặt các gói cần thiết cho việc build
RUN apt-get update && apt-get install -y pkg-config libasound2-dev

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o text-to-voice .

# Giai đoạn runtime
FROM golang:1.21

# Cài đặt mplayer và các thư viện runtime cần thiết
RUN apt-get update && apt-get install -y mplayer libasound2

WORKDIR /app

COPY --from=builder /app/text-to-voice .

# Tạo thư mục audio nếu chưa tồn tại
RUN mkdir -p audio

EXPOSE 5001

CMD ["./text-to-voice"]
