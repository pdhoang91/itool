FROM python:3.8-slim

WORKDIR /app

# Cài đặt các gói hệ thống cần thiết
RUN apt-get update && apt-get install -y \
    libstdc++6 \
    libgcc-s1 \
    sox \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y wget unzip \
    && wget https://github.com/mozilla/DeepSpeech/releases/download/v0.9.3/deepspeech-0.9.3-models.pbmm \
    && wget https://github.com/mozilla/DeepSpeech/releases/download/v0.9.3/deepspeech-0.9.3-models.scorer \
    && wget https://github.com/mozilla/DeepSpeech/releases/download/v0.9.3/deepspeech-0.9.3-models.tar.gz \
    && tar -xvzf deepspeech-0.9.3-models.tar.gz


# Nâng cấp pip lên phiên bản mới nhất
RUN pip install --upgrade pip

# Copy và cài đặt các thư viện Python
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy mã nguồn vào container
COPY . .

# Mở cổng ứng dụng
EXPOSE 5002

# Chạy ứng dụng
CMD ["python", "app.py"]
