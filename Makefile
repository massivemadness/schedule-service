.PHONY: deps
deps:
	go mod download

.PHONY: build-binary
build-binary:
	go build -o ./build/output/main ./cmd/schedule-service

.PHONY: run
run:
	go run ./cmd/schedule-service

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -r ./build/output