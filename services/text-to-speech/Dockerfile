# Sử dụng Python base image
FROM --platform=linux/arm64/v8 python:3.9-slim

# Cài đặt system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libsndfile1 \
    git \
    pkg-config \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Tạo và activate virtual environment
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Nâng cấp pip và cài đặt wheel
RUN pip install --no-cache-dir --upgrade pip wheel setuptools

# Cài đặt dependencies với versions cụ thể theo thứ tự
RUN pip install --no-cache-dir \
    numpy==1.21.6 \
    torch==2.0.0 \
    scipy==1.10.1 \
    Cython==3.0.0 \
    Werkzeug==2.2.3 \
    flask==2.2.3 \
    soundfile==0.12.1

# Cài đặt TTS với version cũ hơn và stable hơn
RUN pip install --no-cache-dir TTS==0.13.0

# Tạo thư mục làm việc
WORKDIR /app

# Copy source code
COPY . .

# Expose port
EXPOSE 5001

# Command để chạy service
CMD ["python", "app.py"]