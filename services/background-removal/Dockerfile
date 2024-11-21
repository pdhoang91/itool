# Sử dụng hình ảnh cơ bản Python 3.9 slim
FROM python:3.9-slim

# Đặt thư mục làm việc
WORKDIR /app

# Cài đặt wget và các phụ thuộc hệ thống khác nếu cần
RUN apt-get update && apt-get install -y wget && rm -rf /var/lib/apt/lists/*

# Sao chép tệp requirements.txt vào thư mục làm việc
COPY requirements.txt .

# Cài đặt các thư viện Python cần thiết
RUN pip install --no-cache-dir -r requirements.txt

# Sao chép toàn bộ mã nguồn vào thư mục làm việc
COPY . .

# Tạo thư mục để lưu log và đảm bảo quyền truy cập
RUN mkdir -p /app/logs
RUN touch /app/logs/app.log
RUN chmod 666 /app/logs/app.log

# Tạo thư mục để lưu model và đảm bảo quyền truy cập
RUN mkdir -p /root/.u2net
RUN chmod 777 /root/.u2net

# Tải model u2net.onnx trước
RUN wget -O /root/.u2net/u2net.onnx https://github.com/danielgatis/rembg/releases/download/v0.0.0/u2net.onnx

# Mở cổng 5003
EXPOSE 5003

# Chạy ứng dụng Flask
CMD ["python", "app.py"]
