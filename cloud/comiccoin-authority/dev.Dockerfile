# The base go-image
FROM golang:1.23

COPY . /go/src/github.com/LuchaComics/monorepo/cloud/comiccoin-authority
WORKDIR /go/src/github.com/LuchaComics/monorepo/cloud/comiccoin-authority

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy only `.go` files, if you want all files to be copied then replace `with `COPY . .` for the code below.
COPY *.go .

# Install our third-party application for hot-reloading capability.
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -polling=true -log-prefix=false -build="go build ." -command="./comiccoin-authority daemon" -directory="./"

# BUILD
# docker build --rm -t comiccoin-authority -f dev.Dockerfile .

# EXECUTE
# docker run -d -p 8000:8000 comiccoin-authority
