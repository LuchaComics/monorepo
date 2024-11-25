# DEVELOPERS NOTE:
# THE PURPOSE OF THIS DOCKERFILE IS TO BUILD THE COMICCOIN-FAUCET EXECUTABLE
# IN A CONTAINER FOR DEVELOPMENT PURPOSES ON YOUR DEVELOPMENT MACHINE. DO
# NOT RUN THIS IN PRODUCTION ENVIRONMENT.

# The base go-image
FROM golang:1.23

COPY . /go/src/github.com/LuchaComics/monorepo/cloud/comiccoin-faucet
WORKDIR /go/src/github.com/LuchaComics/monorepo/cloud/comiccoin-faucet

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy only `.go` files, if you want all files to be copied then replace `with `COPY . .` for the code below.
COPY *.go .

# Install our third-party application for hot-reloading capability.
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -polling=true -log-prefix=false -build="go build ." -command="./comiccoin-faucet daemon" -directory="./"

# BUILD
# docker build --rm -t comiccoin-faucet -f dev.Dockerfile .

# EXECUTE
# docker run -d -p 8000:8000 comiccoin-faucet
