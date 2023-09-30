ENTRYPOINT=cmd/app/main.go
DOC_DIR=internal/api/docs

build:
	go build -o bin/app cmd/app/*.go

run: build
	./bin/app

doc:
	swag fmt
	swag init -g $(ENTRYPOINT) --pd -o $(DOC_DIR)
