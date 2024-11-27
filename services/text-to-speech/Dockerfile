FROM --platform=linux/arm64/v8 python:3.9-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libsndfile1 \
    git \
    pkg-config \
    ffmpeg \
    libsndfile1 \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

RUN pip install --no-cache-dir --upgrade pip wheel setuptools

RUN pip install --no-cache-dir \
    numpy==1.21.6 \
    torch==2.0.0 \
    scipy==1.10.1 \
    Cython==3.0.0 \
    Werkzeug==2.2.3 \
    flask==2.2.3 \
    soundfile==0.12.1 \
    librosa==0.10.0

RUN pip install --no-cache-dir TTS==0.13.0

WORKDIR /app
COPY . .
EXPOSE 5001
CMD ["python", "app.py"]