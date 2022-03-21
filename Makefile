all: test vet fmt lint build

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l
	test -z $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l)

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

build:
	GOOS=linux GOARCH=amd64 go build -o bin/task-runner ./cmd/task-runner
	GOOS=linux GOARCH=amd64 go build -o bin/web-server ./cmd/web-server

run-task-runner:
	go run cmd/task-runner/main.go

run-web-server:
	go run cmd/web-server/main.go