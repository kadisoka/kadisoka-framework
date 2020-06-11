
GOLANG_IMAGE ?= golang:1.14
TESTER_IMAGE ?= citadel-tester

.PHONY: fmt test deps-up

fmt:
	@echo "Formatting files..."
	@docker run --rm \
		-v $(CURDIR):/workspace \
		--entrypoint gofmt \
		$(GOLANG_IMAGE) -w -l -s \
		.

test:
	@echo "Preparing test runner..."
	@docker build -t $(TESTER_IMAGE) -f ./tools/tester.dockerfile . > /dev/null
	@echo "Executing unit tests..."
	@docker run --rm \
		-v $(CURDIR):/workspace \
		--workdir /workspace \
		$(TESTER_IMAGE) test -v ./...

deps-up:
	@echo "Updating all dependencies..."
	@docker run --rm \
		-v $(CURDIR):/workspace \
		--workdir /workspace \
		$(GOLANG_IMAGE) /bin/sh -c "go get -u all && go mod tidy"
