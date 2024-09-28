# Requests

To self-host:
```console
docker compose -f docker-compose.prod.yml up
```

Main page runs on `http://localhost:5000`. 

Webhook page runs on `http://localhost:8080` however it is setup to only handle requests of the form `*.a.b.c`.

On production, I run this with an nginx proxy routing requests internally.
