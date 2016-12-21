VERSION := $(shell cat version.go | grep -Eo "[0-9]+\.[0-9]+\.[0-9]+")

pypihub: ./*.go ./cmd/pypihub/*.go
	go build -o pypihub ./cmd/pypihub/

build/pypihub: ./*.go ./cmd/pypihub/*.go
	mkdir -p build/
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w' -o build/pypihub ./cmd/pypihub/

clean:
	rm -f ./pypihub
	rm -rf ./build

docker_build: build/pypihub
	docker build -t brettlangdon/pypihub .

docker_up: docker_build
		docker-compose up --build

build_release: clean
	# Darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -a -ldflags '-w' -o build/pypihub.${VERSION}.darwin_386 ./cmd/pypihub/
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags '-w' -o build/pypihub.${VERSION}.darwin_amd64 ./cmd/pypihub/
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -ldflags '-w' -o build/pypihub.${VERSION}.linux_386 ./cmd/pypihub/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o build/pypihub.${VERSION}.linux_amd64 ./cmd/pypihub/
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -ldflags '-w' -o build/pypihub.${VERSION}.linux_arm ./cmd/pypihub/
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -ldflags '-w' -o build/pypihub.${VERSION}.windows_386 ./cmd/pypihub/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags '-w' -o build/pypihub.${VERSION}.windows_amd64 ./cmd/pypihub/

.PHONY: build_release clean docker_build docker_up
