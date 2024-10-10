BACKEND_FILES = $(shell find backend -type f -name '*')
FRONTEND_FILES = $(shell find frontend/src frontend/static frontend/*.js* -type f -name '*')

.PHONY: all
all: dist yt-archive rebuild_mpd

dist: $(FRONTEND_FILES)
	cd frontend &&\
	npm install &&\
	npm run build

yt-archive: $(BACKEND_FILES)
	cd backend &&\
	go build -o ..

rebuild_mpd: yt-archive
	cd backend &&\
	go build -o .. ./cmd/rebuild_mpd

.PHONY: standalone
standalone: dist $(BACKEND_FILES)
	cp -r dist backend &&\
	cd backend &&\
	go build -tags standalone -o .. &&\
	rm -r dist

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
	rm rebuild_mpd
	rm -r dist