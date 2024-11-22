FROM python:3.9-slim

WORKDIR /app

# Cài đặt các dependencies cho OpenCV và Tesseract
RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    libtesseract-dev \
    tesseract-ocr-vie \
    libgl1-mesa-glx \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

# Cài đặt numpy trước
RUN pip install numpy==1.24.3

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 5006

CMD ["python", "app.py"]