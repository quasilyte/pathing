GOPATH_DIR=`go env GOPATH`

test:
	go test -count 2 ./...
	go test -run=__ -bench=. ./...
	cd _bench && go test -v -bench=. .
	@echo "everything is OK"

ci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.54.2
	$(GOPATH_DIR)/bin/golangci-lint run ./...
	@echo "everything is OK"

lint:
	golangci-lint run ./...
	@echo "everything is OK"

.PHONY: ci-lint lint test
