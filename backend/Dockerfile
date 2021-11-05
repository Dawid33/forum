# Build the static files and executable backend
FROM alpine:edge

RUN apk update && apk upgrade

# Install dependencies
RUN apk add --no-cache go

# Copy all files to docker image
RUN mkdir -p /backend/forum-templates && mkdir -p /website/main
COPY main /backend/main/
COPY forum-templates /backend/forum-templates/

WORKDIR /backend/main
RUN go build -o backend && mv backend ..
WORKDIR /backend
CMD ./backend