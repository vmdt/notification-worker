GOPATH:=$(shell go env GOPATH)

.PHONY: run
run:
	@echo "Running Go application..."
	@cd cmd && go run main.go