build:
	go build -o ./bin/main ./cmd/blocd/
	sudo cp ./bin/main /usr/local/bin/blocd