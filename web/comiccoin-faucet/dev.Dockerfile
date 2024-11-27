# build environment
FROM node:14-alpine as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH

# Get environment variables.
ARG REACT_APP_API_HOST
ARG REACT_APP_API_DOMAIN
ARG REACT_APP_API_PROTOCOL
ARG REACT_APP_WWW_DOMAIN
ARG REACT_APP_WWW_PROTOCOL

COPY ./package.json ./
COPY ./package-lock.json ./

RUN npm install

COPY . ./

# Apply environment variables with the build.
RUN REACT_APP_API_HOST=${REACT_APP_API_HOST} \
  REACT_APP_API_DOMAIN=${REACT_APP_API_DOMAIN} \
  REACT_APP_API_PROTOCOL=${REACT_APP_API_PROTOCOL} \
  REACT_APP_WWW_DOMAIN=${REACT_APP_WWW_DOMAIN} \
  REACT_APP_WWW_PROTOCOL=${REACT_APP_WWW_PROTOCOL}

# pr
EXPOSE 3000
CMD ["npm", "run", "start"]
