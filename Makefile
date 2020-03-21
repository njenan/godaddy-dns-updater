test:
	go test -v ./...

#generate:
	# statik -src=./ -dest=./ -p=cmd

build: #generate
	go build -o godaddy-dns-updater -ldflags "-X main.Version=$(shell cat version) -X main.License=GPL-2.0" .

dockerize: build
	docker build . -t docker.nathanjenan.me/njenan/godaddy-dns-updater:latest --add-host=athens.nathanjenan.me:192.168.0.200
	docker build . -t docker.nathanjenan.me/njenan/godaddy-dns-updater:$(shell git rev-parse --short HEAD) --add-host=athens.nathanjenan.me:192.168.0.200
	docker push docker.nathanjenan.me/njenan/godaddy-dns-updater:latest
	docker push docker.nathanjenan.me/njenan/godaddy-dns-updater:$(shell git rev-parse --short HEAD)

