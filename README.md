# Requests

To self-host:
```console
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

Main page runs on `http://localhost:5000`. 

Webhook page runs on `http://localhost:8080` however it is setup to only handle requests of the form `*.a.b.c`.

On production, I run this with an nginx proxy routing requests internally.


## TODO:
- [ ] Add API routes to customize the Bin response
- [ ] Improve the enviornment to support extracting the bin dynamically based on user criteria
- [ ] Look into further redis optimizations?
- [ ] Maybe use Redis Streams instead of PubSub?
