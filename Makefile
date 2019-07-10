.PHONY: build
build:
	go build -o imgcontent

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofmt -w -s ./

.PHONY: lint
lint:
	golint --set_exit_status ./...

