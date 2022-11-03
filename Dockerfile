FROM python:3.10.6-slim as app
WORKDIR /app
COPY requirements.txt .
RUN apt-get update &&\
    apt-get install --no-install-recommends -y build-essential python3-dev &&\
    apt-get purge -y build-essential &&\
    rm -rf /var/lib/apt/lists/* &&\
    pip install -r requirements.txt --no-cache-dir
COPY . .
CMD ["python", "main.py"]
