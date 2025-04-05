COMMON_BACKEND_FILES := $(wildcard backend/go.* backend/**/*.go)

BACKEND_FILES := $(wildcard backend/cmd/yt-archive/*.go) $(COMMON_BACKEND_FILES)
CLI_FILES := $(wildcard backend/cmd/yt-archive-cli/*.go) $(COMMON_BACKEND_FILES)

FRONTEND_FILES := $(wildcard frontend/src/* frontend/src/**/* frontend/static/**/* frontend/*.js*)

SQL_FILES := $(wildcard backend/db/sql/*)

all: dist yt-archive yt-archive-cli

dist: $(FRONTEND_FILES)
	cd frontend &&\
	npm install &&\
	npm run build

sql: $(SQL_FILES)
	cd backend &&\
	sqlc generate

yt-archive: $(BACKEND_FILES)
	cd backend &&\
	go build -o .. ./cmd/yt-archive

standalone: dist $(BACKEND_FILES)
	cp -r dist backend/cmd/yt-archive &&\
	cd backend &&\
	go build -tags standalone -o .. ./cmd/yt-archive &&\
	rm -r cmd/yt-archive/dist

yt-archive-cli: $(CLI_FILES)
	cd backend &&\
	go build -o .. ./cmd/yt-archive-cli

start: all
	./yt-archive

dev: yt-archive
	YT_ARCHIVE_ADDR=localhost:8080 ./yt-archive & \
	trap 'kill $$(pgrep yt-archive) && exit 0' SIGINT SIGTERM && \
	(cd frontend && npm run dev) 

clean:
	rm yt-archive yt-archive-cli
	rm -r dist

reset:
	rm -r thumbnails videos database

.PHONY: all standalone start dev clean reset