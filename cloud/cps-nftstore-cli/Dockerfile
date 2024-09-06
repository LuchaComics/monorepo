# node:alpine will be our base image to create this image
FROM node:18-alpine

# Set the /app directory as working directory
WORKDIR /app

# Install ganache-cli globally (Note: https://github.com/clearmatics/ganache-cli)
RUN npm install -g ganache-cli

# Set the default command for the image
CMD ["ganache-cli", "-v=true", "-h", "0.0.0.0"]

# To startup and execute:
# docker build -t ganache .
# docker run -p 8545:8545 ganache
