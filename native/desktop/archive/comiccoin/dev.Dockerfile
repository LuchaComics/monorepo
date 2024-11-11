# Developers Note:
# 1. The purpose of this Dockerfile is to provide you with an environment for development, not production deployment!
# 2. Before beginning please setup a `identitykey`.
# 3. If you are running in `host mode` then do not set `COMICCOIN_BOOTSTRAP_PEERS` and leave as empty string.
# 4. If you are running in `dial mode` then set `COMICCOIN_BOOTSTRAP_PEERS`.

# Environment variables we will pass into this docker file from our docker-compose (or command line if you prefer to run without docker-compose).
ARG COMICCOIN_DATADIR
ARG COMICCOIN_IDENTITYKEY_ID
ARG COMICCOIN_LISTEN_HTTP_ADDRESS
ARG COMICCOIN_LISTEN_P2P_PORT
ARG COMICCOIN_BOOTSTRAP_PEERS

# The base go-image
FROM golang:1.23

COPY . /go/src/github.com/LuchaComics/monorepo/native/desktop/comiccoin
WORKDIR /go/src/github.com/LuchaComics/monorepo/native/desktop/comiccoin

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy only `.go` files, if you want all files to be copied then replace `with `COPY . .` for the code below.
COPY *.go .

# Install our third-party application for hot-reloading capability.
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -polling=true -log-prefix=false -build="go build ." -directory="./" -command="./comiccoin daemon --datadir=${COMICCOIN_DATADIR} --identitykey-id=${COMICCOIN_IDENTITYKEY_ID} --listen-http-address=${COMICCOIN_LISTEN_HTTP_ADDRESS} --listen-p2p-port=${COMICCOIN_LISTEN_P2P_PORT} --bootstrap-peers=${COMICCOIN_BOOTSTRAP_PEERS}"

# Developer Notes:
#
#    BUILD
#    docker build --rm -t comiccoin -f dev.Dockerfile .
#
#    EXECUTE
#    docker run -d -p 8000:8000 comiccoin
