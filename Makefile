test:
	go test -v ./...

build:
	go build -o godaddy-dns-updater -ldflags "-X main.Version=$(shell cat version)" .

version:
	git describe --tags --abbrev=0 | tee version

dockerize: version build
	docker build . -t docker.nathanjenan.me/njenan/godaddy-dns-updater:latest
	docker push docker.nathanjenan.me/njenan/godaddy-dns-updater:latest
