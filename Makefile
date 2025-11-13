.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -cover ./...

.PHONY: build
build:
	go build -o pudd .

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -f pudd
