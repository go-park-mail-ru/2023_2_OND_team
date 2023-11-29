.PHONY: build run test test_with_coverage cleantest retest doc generate cover_all currcover
.PHONY: build_auth  build_realtime build_messenger

ENTRYPOINT=cmd/app/main.go
DOC_DIR=./docs
COV_OUT=coverage.out
COV_HTML=coverage.html
CURRCOVER=github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1

build:
	go build -o bin/app cmd/app/*.go

build_auth:
	go build -o bin/auth cmd/auth/*.go

build_realtime:
	go build -o bin/realtime cmd/realtime/*.go

build_messenger:
	go build -o bin/messenger cmd/messenger/*.go

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

generate:
	go generate ./...

cover_all:
	go test -coverpkg=./... -coverprofile=cover ./...
	cat cover | grep -v "mock" | grep -v  "easyjson" | grep -v "proto" | grep -v "ramrepo" > cover.out
	go tool cover -func=cover.out

currcover:
	go test -cover -v -coverprofile=cover.out ${CURRCOVER}
	go tool cover -html=cover.out -o cover.html

