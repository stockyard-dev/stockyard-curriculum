build:
	CGO_ENABLED=0 go build -o curriculum ./cmd/curriculum/

run: build
	./curriculum

test:
	go test ./...

clean:
	rm -f curriculum

.PHONY: build run test clean
