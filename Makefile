pypihub: ./*.go
	go build -o pypihub ./cmd/pypihub/

clean:
	rm -f ./pypihub

run:
	go run ./cmd/pypihub/main.go

.PHONY: clean
