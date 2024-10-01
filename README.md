# Requests

## Recommended self-host steps:

Modify the docker-compose-prod with any env variables.
```console
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

Main page runs on `http://localhost:5000`. 

Webhook page runs on `http://localhost:8080`. 

## Local Development
- TODO: Make a devcontainer
- The provided docker-compose.yml runs the local build version of the site. 
- Run the go backend with `go run main.go`
- Run the flask server in dev mode with reload, use the build script to write front end changes to the static directory of the flask server.

## Enviornment Variables
- **REDIS_URL**
    - Specify the endpoint for the redis url in the form host:port
    - Default: redis:6379 using the internal docker network
- **BIN_REGEX**
    - Specify how to extract the bin id from the url
    - Default: `^(?:[a-zA-Z0-9-_.]+?)(?:\:\d+)?/(?<bin>\w*)(?:/.*)?$`
        - This will extract the Bin as the first path from the URL.
        - Examples: 
            - localhost:8080/bin/asdf
            - r.ahh.bet/bin
            - r.ahh.bet/bin?asdf


## TODO:
- Add API routes and UI to customise the response content
- I should prob host it somewhere...

