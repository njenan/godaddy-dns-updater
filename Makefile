test:
	go test -v ./...

build: version
	go build -o godaddy-dns-updater -ldflags "-X main.Version=$(shell cat version)" .

install: version
	go install -ldflags "-X main.Version=$(shell cat version)"

dockerize: build
	docker build . -t docker.nathanjenan.me/njenan/godaddy-dns-updater:latest
	docker build . -t docker.nathanjenan.me/njenan/godaddy-dns-updater:$(shell git rev-parse --short HEAD)
	docker push docker.nathanjenan.me/njenan/godaddy-dns-updater:latest
	docker push docker.nathanjenan.me/njenan/godaddy-dns-updater:$(shell git rev-parse --short HEAD)

