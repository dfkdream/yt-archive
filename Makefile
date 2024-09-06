all: frontend backend

frontend: dist

backend: yt-archive

dist: 
	cd frontend &&\
	npm install &&\
	npm run build

yt-archive: 
	cd backend &&\
	go build -o ..

clean:
	rm yt-archive
	rm -r dist