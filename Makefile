ENTRYPOINT=cmd/app/main.go
DOC_DIR=./docs
COV_OUT=coverage.out
COV_HTML=coverage.html

build:
	go build -o bin/app cmd/app/*.go

run: build
	./bin/app

test:
	go test ./... -race -covermode=atomic -coverpkg ./... -coverprofile=$(COV_OUT)

test_with_coverage: test
	go tool cover -html $(COV_OUT) -o $(COV_HTML)

cleantest:
	rm coverage*

retest:
	make cleantest
	make test

doc:
	swag fmt
	swag init -g $(ENTRYPOINT) --pd -o $(DOC_DIR)
