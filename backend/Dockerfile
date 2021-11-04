# Build the static files and executable backend
FROM alpine AS builder

RUN apk update

# Install dependencies
RUN apk add --no-cache go

# Copy all files to docker image
RUN mkdir -p /backend/forum-templates && mkdir -p /website/main
COPY main /backend/main/
COPY forum-templates /backend/forum-templates/

WORKDIR /backend/main
RUN go build && mv backend ..
CMD /backend/backend