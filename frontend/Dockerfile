# Build the static files and executable backend
FROM alpine as builder

RUN apk update

# Install dependencies
RUN apk add --no-cache nodejs npm curl && mkdir /website

# Copy all files to docker image
RUN mkdir -p /website/src && mkdir -p /website/public
ADD ./src/ /website/src
ADD ./public/ /website/public
COPY ["package.json", "package-lock.json", "tsconfig.json", "/website/"]

# Build frontend
WORKDIR /website
RUN npm i react-scripts typescript && npm run build

FROM nginx
COPY --from=builder /website/build /usr/share/nginx/html/