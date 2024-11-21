FROM python:3.9

WORKDIR /app

# Install necessary system packages
RUN apt-get update && apt-get install -y \
    libstdc++6 \
    libgcc-s1 \
    sox \
    && rm -rf /var/lib/apt/lists/*

# Upgrade pip to the latest version
RUN pip install --upgrade pip

# Copy and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy your application code
COPY . .

# Expose the application port
EXPOSE 5002

# Run the application
CMD ["python", "app.py"]
