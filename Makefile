build:
	go build -o ./bin/ec2update

run: build
	./bin/ec2update

test:
	@go test -v ./...
