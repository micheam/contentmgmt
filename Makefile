APPNAME=imgcontent

$(APPNAME): clean
	mkdir ./bin && go build -o ./bin/imgcontent

.PHONY: clean
clean:
	rm -rf ./bin/

.PHONY: install
install:
	go install

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofmt -w -s ./

.PHONY: lint
lint:
	golint --set_exit_status ./...

