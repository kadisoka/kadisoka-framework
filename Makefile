
GOLANG_IMAGE ?= golang:1.18

.PHONY: deps-up
deps-up:
	@echo "Updating all dependencies..."
	@docker run --rm \
		-v $(CURDIR):/workspace \
		--workdir /workspace \
		$(GOLANG_IMAGE) /bin/sh -c "go get -u all && go mod tidy"
