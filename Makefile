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

clean:
	rm yt-archive
	rm -r dist