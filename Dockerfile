FROM node AS frontend

WORKDIR /yt-archive

COPY ./frontend ./frontend

COPY ./Makefile .

RUN make dist

FROM golang:alpine AS backend

RUN apk update

RUN apk add build-base

WORKDIR /yt-archive

COPY ./backend ./backend

COPY ./Makefile .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make yt-archive

FROM alpine

RUN apk update

RUN apk add yt-dlp ffmpeg

WORKDIR /yt-archive

COPY --from=frontend /yt-archive/dist ./dist

COPY --from=backend /yt-archive/yt-archive ./yt-archive

ENTRYPOINT ["./yt-archive"]