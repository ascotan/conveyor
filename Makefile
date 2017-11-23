all: build

init:
	dep ensure

build: init
	mkdir -p build
	go vet ./...
	go fmt ./...
	go build -o build/conveyor -race -v cmd/conveyor/conveyor.go

clean:
	rm -rf build vendor

serve:
	./build/conveyor


.PHONY: build clean