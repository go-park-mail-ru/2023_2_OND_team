.PHONY: build run test test_with_coverage cleantest retest doc gen
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

gen:
	mockgen -source=internal/pkg/repository/board/repo.go -source=internal/pkg/repository/user/repo.go \
	-destination=internal/pkg/repository/board/mock/mock_repo.go -destination=internal/pkg/repository/user/mock/mock_repo.go