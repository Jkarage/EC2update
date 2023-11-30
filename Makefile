build:
	@go build -o ./bin/amiupdate

test:
	@go test -v ./...

run: build
	./bin/amiupdate