# build stage
FROM golang:1.20 AS build

ENV GO111MODULE=on

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=build /app/main .

# Define your environment variables here
ENV PORT=8080 \
    LOG_FILE=/var/log/app/golang-app.log \
    MOVIE_SOURCE_FOLDER=/app/movies-source \
    MOVIE_TARGET_FOLDER=/app/movies-target \
    TV_SOURCE_FOLDER=/app/tvshows-source \
    TV_TARGET_FOLDER=/app/tvshows-target \
    TMDB_API_KEY=""

# Expose the port on which the application will listen
EXPOSE $PORT

VOLUME /var/log/app \
         /app/movies-source \
         /app/movies-target \
         /app/tvshows-source \
         /app/tvshows-target

# Start the application
CMD ["/app/main","-serve"]
