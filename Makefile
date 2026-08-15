.PHONY: default test lint vet fix fix-do generate clean-generated tidy

default: lint test

lint: vet fix

test:
	go test -cover -race ./...

vet:
	go vet ./...

fix:
	go fix -diff ./...

fix-do:
	go fix ./...

generate:
	go generate ./...

clean-generated:
	find . -name "*_generated.go" -type f -delete

tidy:
	go mod tidy -v
