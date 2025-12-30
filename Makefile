build:
	go build -o ./bin/main ./cmd/blocd/main.go
	sudo cp ./bin/main /usr/local/bin/blocd