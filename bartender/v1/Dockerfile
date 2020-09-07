FROM python:3-slim AS builder
ENV PYTHONUNBUFFERED=1

WORKDIR /test

# TODO(hvl): unit tests


# Next stage
FROM python:3-slim
ENV PYTHONUNBUFFERED=1

WORKDIR /app

# Pre-install local library
COPY ./var/ ./var/
RUN pip install var/.

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY src .

ENTRYPOINT ["python", "server.py"]
