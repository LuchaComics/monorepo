# node:alpine will be our base image to create this image
FROM node:18-alpine

# Set the /app directory as working directory
WORKDIR /app

# Install ganache-cli globally (Note: https://github.com/clearmatics/ganache-cli)
RUN npm install -g ganache-cli

# DEVELOPERS NOTES:
# There are two options we can use for testing purposes in our development environment:
# Option 1: Use this option if you would like to generate a new Ethereum test blockchain every time this docker loads up, e.i. data is temporary.
# Option 2: Use this option to make the blockchain changes permanent so every time you restart this docker container then the data is persistent.

# (Option 1) Set the default command for the image to load up our own testnet.
# CMD ["ganache-cli", "-v=true", "-h", "0.0.0.0"]

# (Option 2) Set the default command for the image to load up our own testnet.
CMD ["ganache-cli", "-v=true", "-h", "0.0.0.0", "--db=/ganache_data", "--mnemonic", "'ski crystal tell aware quote dash live turn awkward seat zone skill'"]

# To startup and execute:
# docker build -t ganache-cli -f ganache-cli.Dockerfile .
# docker run --name GanacheCLI -p 8545:8545 ganache-cli
