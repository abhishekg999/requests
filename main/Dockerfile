FROM python:3.10-slim

LABEL org.opencontainers.image.source="https://github.com/abhishekg999/requests"

WORKDIR /app

COPY backend/requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt

ENV GUNICORN_CMD_ARGS="--bind=0.0.0.0"
COPY backend/. .

EXPOSE 5000

CMD [ "gunicorn", "-b", ":5000", "--threads", "100", "app:app" ]