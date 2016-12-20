pypihub: ./*.go ./cmd/pypihub/*.go
	go build -o pypihub ./cmd/pypihub/

build/pypihub: ./*.go ./cmd/pypihub/*.go
	mkdir -p build/
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o build/pypihub ./cmd/pypihub/

clean:
	rm -f ./pypihub
	rm -rf ./build

run:
	go run ./cmd/pypihub/main.go

docker_build: build/pypihub
	docker build -t pypihub .

docker_up: docker_build
		docker-compose up --build

.PHONY: clean docker_build run docker)up
