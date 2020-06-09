PKG_PATH = github.com/citadelium/foundation/pkg
GOLANG_IMAGE ?= golang:1.14

.PHONY: fmt deps-up

fmt:
	@echo "Formatting files..."
	@docker run --rm \
		-v $(CURDIR):/go \
		--entrypoint gofmt \
		$(GOLANG_IMAGE) -w -l -s \
		.

deps-up:
	@echo "Updating all dependencies..."
	@docker run --rm \
		-v $(CURDIR):/$(PKG_PATH) \
		--workdir /$(PKG_PATH) \
		$(GOLANG_IMAGE) /bin/sh -c "go get -u all && go mod tidy"
