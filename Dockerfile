# Build the static files and executable backend
FROM alpine AS builder

RUN apk update

# Install dependencies
RUN apk add --no-cache nodejs npm curl go

# Copy all files to docker image
RUN mkdir -p /website/public/forum && mkdir -p /website/public/console && mkdir -p /website/frontend && mkdir -p /website/main
WORKDIR /website
COPY ./frontend ./frontend
COPY ./main ./main
COPY ./public/forum ./public/forum/

# Build frontend
WORKDIR /website/frontend
RUN npm i react-scripts typescript && npm run build
RUN cp -r ./build/* ../public/console
# This is temporary until I get the backend to work with environment variables
WORKDIR /website/main
RUN go build && mv backend .. && ln -sf /website/public/console /website/public/console-dev

WORKDIR /website
CMD /website/backend



# Run backend


