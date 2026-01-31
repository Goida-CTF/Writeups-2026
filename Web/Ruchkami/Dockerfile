FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    chromium \
    chromium-driver \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

# Добавляем CA сертификат в системные доверенные
RUN cp /app/certs/ca.crt /usr/local/share/ca-certificates/ruchkami-ca.crt && \
    update-ca-certificates

EXPOSE 31337

CMD ["python", "app.py"]