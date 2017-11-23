build:
	go build -o build/conveyor -race -v cmd/conveyor/conveyor.go

clean:
	rm -f build/conveyor


.PHONY: build clean