
GOLANG_IMAGE ?= golang:1.16
TESTER_IMAGE ?= kadisoka-foundation-tester

.PHONY: fmt
fmt:
	@echo "Formatting files..."
	@docker run --rm \
		-v $(CURDIR):/go \
		--entrypoint gofmt \
		$(GOLANG_IMAGE) -w -l -s \
		./pkg

# TODO: move this to root?
.PHONY: test
test:
	@echo "Preparing test runner..."
	@docker build -t $(TESTER_IMAGE) -f ./tools/tester.dockerfile . > /dev/null
	@echo "Executing unit tests..."
	@docker run --rm \
		-v $(CURDIR):/workspace \
		--workdir /workspace \
		$(TESTER_IMAGE) test -v ./...
