# Build the static files and executable backend
FROM alpine

RUN apk update

# Install dependencies
RUN apk add --no-cache nodejs npm curl && mkdir /website

# Copy all files to docker image
COPY ./* ./website/

# Build frontend
WORKDIR /website
RUN npm i react-scripts typescript
CMD npm run start