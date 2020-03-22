test:
	go test -v ./...

dockerize:
	go build -o godaddy-dns-updater .
	docker build . -t docker.nathanjenan.me/njenan/godaddy-dns-updater:latest
	docker push docker.nathanjenan.me/njenan/godaddy-dns-updater:latest
