# Build environment.
FROM node:20-alpine AS builder
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH

# Get environment variables from the build.
ARG REACT_APP_API_HOST
ARG REACT_APP_API_DOMAIN
ARG REACT_APP_API_PROTOCOL
ARG REACT_APP_WWW_DOMAIN
ARG REACT_APP_WWW_PROTOCOL

# Set environment variables in this image when built.
# Apply environment variables with the build.
RUN REACT_APP_API_HOST=${REACT_APP_API_HOST} \
  REACT_APP_API_DOMAIN=${REACT_APP_API_DOMAIN} \
  REACT_APP_API_PROTOCOL=${REACT_APP_API_PROTOCOL} \
  REACT_APP_WWW_DOMAIN=${REACT_APP_WWW_DOMAIN} \
  REACT_APP_WWW_PROTOCOL=${REACT_APP_WWW_PROTOCOL}

# Verify environment variables work.
RUN echo "$REACT_APP_API_HOST"
RUN echo "$REACT_APP_API_DOMAIN"
RUN echo "$REACT_APP_API_PROTOCOL"
RUN echo "$REACT_APP_WWW_DOMAIN"
RUN echo "$REACT_APP_WWW_PROTOCOL"

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package*.json /app/

RUN npm install --production

# Bundle app source
COPY ./ /app/

RUN npm run build

# Copy into minimalist Nginx server.
# # SPECIAL THANKS: https://rsbh.dev/blog/dockerize-react-app
FROM nginx:1.24-alpine AS server
COPY --from=builder /app/build/ /usr/share/nginx/html
COPY --from=builder /app/nginx/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

# run nginx with global directives and daemon off
CMD ["nginx", "-g", "daemon off;"]

### BUILD
# docker build -f prod.Dockerfile -t rodolfossp/comiccoin-faucet:prod-latest --platform linux/amd64 --build-arg REACT_APP_API_HOST=https://cpsdata.ca --build-arg REACT_APP_API_DOMAIN=cpsdata.ca --build-arg REACT_APP_API_PROTOCOL=https --build-arg REACT_APP_WWW_DOMAIN=cpsapp.ca --build-arg REACT_APP_WWW_PROTOCOL=https .

### TAG
# docker tag rodolfossp/comiccoin-faucet:prod-latest rodolfossp/comiccoin-faucet:prod-latest

### UPLOAD
# docker push rodolfossp/comiccoin-faucet:prod-latest
