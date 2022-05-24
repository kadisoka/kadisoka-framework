
LINTER_IMAGE ?= kadisoka-linter
TESTER_IMAGE ?= kadisoka-tester
GOLANG_IMAGE ?= golang:1.18
POSTGRES_IMAGE ?= postgres:10.7-alpine

DOCKER_NETWORK ?= iam_default

.PHONY: run
run:
	@docker-compose up --build

.PHONY: fmt
fmt:
	@echo "Formatting files..."
	@docker run --rm \
		-v $(CURDIR):/go \
		--entrypoint gofmt \
		$(GOLANG_IMAGE) -w -l -s \
		./pkg ./iam-server ./examples

.PHONY: lint
lint:
	@echo "Preparing linter..."
	@docker build -t $(LINTER_IMAGE) -f ./tools/linter.dockerfile . > /dev/null
	@echo "Running lint..."
	@docker run --rm \
		-v $(CURDIR):/workspace \
		--workdir /workspace \
		$(LINTER_IMAGE) \
		./pkg/... ./iam-server/... ./examples/...

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

# audit: https://github.com/securego/gosec

IAM_DB_PASSWORD ?=
IAM_DB_URL ?= postgres://iam:$(IAM_DB_PASSWORD)@iam-db/iam?sslmode=disable

# TODO: build this into the executable. migrations should be baked into the executable.
.PHONY: _init_iam_db
_init_iam_db:
	docker run --rm -v $(CURDIR)/pkg/iamserver/migrations:/migrations \
		--entrypoint=psql --network=$(DOCKER_NETWORK) \
		$(POSTGRES_IMAGE) \
		-f ./migrations/000.00-iam.sql \
		"$(IAM_DB_URL)"
