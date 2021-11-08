# Build the static files and executable backend
FROM alpine:edge

RUN apk update && apk upgrade

# Install dependencies
RUN apk add --no-cache go

# Copy all files to docker image
RUN mkdir -p /backend/public && mkdir -p /website/main
COPY main /backend/main/
COPY public /backend/public/

WORKDIR /backend/main
RUN go mod tidy
RUN go build -o backend && mv backend ..
WORKDIR /backend
CMD ./backend