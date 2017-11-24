all: build serve

build: client server

client:
	npm install -g webpack
	cd web; npm install
	cd web; webpack
	cp web/index.html build/

server:
	dep ensure
	mkdir -p build
	go vet ./...
	go fmt ./...
	go build -o build/conveyor -race -v cmd/conveyor/conveyor.go

clean:
	rm -rf build vendor logs web/node_modules

serve:
	cd build; ./conveyor


.PHONY: build clean web serve