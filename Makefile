BACKEND_FILES = $(shell find backend -type f -name '*')
FRONTEND_FILES = $(shell find frontend/src frontend/static frontend/*.js* -type f -name '*')

all: dist yt-archive

dist: $(FRONTEND_FILES)
	cd frontend &&\
	npm install &&\
	npm run build

yt-archive: $(BACKEND_FILES)
	cd backend &&\
	go build -o ..

start: all
	./yt-archive

dev: yt-archive
	YT_ARCHIVE_ADDR=localhost:8080 ./yt-archive &
	(cd frontend && npm run dev)

clean:
	rm yt-archive
	rm -r dist