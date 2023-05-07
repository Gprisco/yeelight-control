run:
	go run main.go

build:
	go build main.go -o bin/yeelightctl

tests:
	go test ./...