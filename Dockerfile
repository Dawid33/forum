# Build the static files and executable backend
FROM alpine AS builder

RUN apk update

# Install dependencies
RUN apk add --no-cache nodejs npm curl go

# Copy all files to docker image
RUN mkdir -p /website/public/forum-templates && mkdir -p /website/public/console && mkdir -p /website/frontend && mkdir -p /website/main
WORKDIR /website
COPY ./frontend ./frontend
COPY ./main ./main
COPY ./public/forum-templates ./public/forum-templates/

# Build frontend
WORKDIR /website/frontend
RUN npm i react-scripts typescript && npm run build-release

# Build backend
WORKDIR /website/main
RUN go build && mv backend ..

WORKDIR /website
CMD /website/backend



# Run backend


