COMMON_BACKEND_FILES := $(wildcard backend/go.* backend/**/*.go)

BACKEND_FILES := $(wildcard backend/cmd/yt-archive/*.go) $(COMMON_BACKEND_FILES)

FRONTEND_FILES := $(wildcard frontend/src/**/* frontend/static/**/* frontend/*.js*)

.PHONY: all
all: dist yt-archive

dist: $(FRONTEND_FILES)
	cd frontend &&\
	npm install &&\
	npm run build

yt-archive: $(BACKEND_FILES)
	cd backend &&\
	go build -o .. ./cmd/yt-archive

.PHONY: standalone
standalone: dist $(BACKEND_FILES)
	cp -r dist backend/cmd/yt-archive &&\
	cd backend &&\
	go build -tags standalone -o .. ./cmd/yt-archive &&\
	rm -r cmd/yt-archive/dist

.PHONY: start
start: all
	./yt-archive

.PHONY: dev
dev: yt-archive
	YT_ARCHIVE_ADDR=localhost:8080 ./yt-archive & \
	trap 'kill $$(pgrep yt-archive) && exit 0' SIGINT SIGTERM && \
	(cd frontend && npm run dev) 

.PHONY: clean
clean:
	rm yt-archive
	rm -r dist