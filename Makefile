build:
	go build -o bin/app cmd/app/*.go

run: build
	./bin/app
