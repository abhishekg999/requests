from __future__ import annotations
import redis
import os
import nanoid

REDIS_URL = os.getenv("REDIS_URL", "localhost:6379")
HOST, PORT = REDIS_URL.split(":")

BIN_PREFIX = "bin:"
CHANNEL_PREFIX = "channel:"

DEFAULT_HTTP_STATUS_CODE = "200"
DEFAULT_HTTP_BODY = "Hi!"
DEFAULT_HTTP_HEADERS = '{"Content-Type": "text/plain"}'

client = redis.Redis(HOST, int(PORT))

def allocateDefaultBin() -> str:
    bin = nanoid.generate("abcdefghijklmnopqrstuvwxyz", 14)
    client.hset(
        BIN_PREFIX + bin,
        mapping={
            "body": DEFAULT_HTTP_BODY,
            "status_code": DEFAULT_HTTP_STATUS_CODE,
            "headers": DEFAULT_HTTP_HEADERS,
        },
    )
    client.expire(BIN_PREFIX + bin, 86400000)
    return bin


def getResponses(bin: str):
    pubsub = client.pubsub()
    pubsub.subscribe(CHANNEL_PREFIX + bin)
    return pubsub


def updateResponse(bin: str, response: dict[str, str]):
    client.hset(BIN_PREFIX + bin, mapping=response)
