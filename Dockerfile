# build stage
FROM golang:1.20 AS build

ENV GO111MODULE=on

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -x -ldflags "-s -w" -o main .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates ffmpeg

WORKDIR /app/
COPY --from=build /app/main .

# Define your environment variables here
ENV GOMEMLIMIT=200MiB \
    TZ=Europe/Paris \
    PORT=8080 \
    LOG_FILE=/app/logs/golang-app.log \
    MOVIE_SOURCE_FOLDER=/app/movies-source \
    MOVIE_TARGET_FOLDER=/app/media-target/movies \
    TV_SOURCE_FOLDER=/app/tvshows-source \
    TV_TARGET_FOLDER=/app/media-target/tv-shows \
    TMDB_API_KEY="" \
    DB_SYNC=true \
    DB_HOST=127.0.0.1 \
    DB_PORT=5432 \
    DB_USER=bingemate \
    DB_PASSWORD=bingemate \
    DB_NAME=bingemate \
    SCAN_CRON="*/15 * * * *"

# Expose the port on which the application will listen
EXPOSE $PORT

VOLUME /var/logs/app \
         /app/movies-source \
         /app/media-target \
         /app/tvshows-source \

USER 1000:100

# Start the application
CMD ["/app/main","-serve"]
