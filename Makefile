# Based on anchore syft and grype repos. Also uses tools made by anchore

BIN := openshift4_mirror
TEMP_DIR := ./.tmp

# Tool versions #################################
GOLANGCILINT_VERSION := v1.52.2
GOSIMPORTS_VERSION := v0.3.8
GORELEASER_VERSION := v1.17.2
OPENSHIFT_VERSION := 4.12.12

# Command templates #################################
GORELEASE_CMD := curl -sfL https://goreleaser.com/static/run | VERSION=$(GORELEASER_VERSION) sh -s --
LINT_CMD := $(TEMP_DIR)/golangci-lint run --tests=false
GOIMPORTS_CMD := $(TEMP_DIR)/gosimports
RELEASE_CMD := $(GORELEASE_CMD) release --clean
SNAPSHOT_CMD := $(RELEASE_CMD) --skip-publish --skip-sign --snapshot

## Build variables #################################
DIST_DIR := ./dist
SNAPSHOT_DIR := ./snapshot

## Bootstrapping targets #################################

.PHONY: bootstrap
bootstrap: $(TEMP_DIR) bootstrap-go bootstrap-tools ## Download and install all tooling dependencies (+ prep tooling in the ./tmp dir)
	$(call title,Bootstrapping dependencies)

.PHONY: bootstrap-tools
bootstrap-tools: $(TEMP_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TEMP_DIR)/ $(GOLANGCILINT_VERSION)
	GOBIN="$(realpath $(TEMP_DIR))" go install github.com/rinchsan/gosimports/cmd/gosimports@$(GOSIMPORTS_VERSION)

.PHONY: bootstrap-go
bootstrap-go:
	go mod download

$(TEMP_DIR):
	mkdir -p $(TEMP_DIR)

## Linting and formating #################################

.PHONY: lint
lint:  ## Run gofmt + golangci lint checks
	$(call title,Running linters)
	# ensure there are no go fmt differences
	@printf "files with gofmt issues: [$(shell gofmt -l -s .)]\n"
	@test -z "$(shell gofmt -l -s .)"

	# run all golangci-lint rules
	$(LINT_CMD)
	@[ -z "$(shell $(GOIMPORTS_CMD) -d .)" ] || (echo "goimports needs to be fixed" && false)

	# go tooling does not play well with certain filename characters, ensure the common cases don't result in future "go get" failures
	$(eval MALFORMED_FILENAMES := $(shell find . | grep -e ':'))
	@bash -c "[[ '$(MALFORMED_FILENAMES)' == '' ]] || (printf '\nfound unsupported filename characters:\n$(MALFORMED_FILENAMES)\n\n' && false)"

.PHONY: format
format:  ## Auto-format all source code
	$(call title,Running formatters)
	gofmt -w -s .
	$(GOIMPORTS_CMD) -w .
	go mod tidy

.PHONY: lint-fix
lint-fix: format  ## Auto-format all source code + run golangci lint fixers
	$(call title,Running lint fixers)
	$(LINT_CMD) --fix

## Build-related targets #################################

.PHONY: build
build: $(SNAPSHOT_DIR)  ## Build release snapshot binaries and packages

$(SNAPSHOT_DIR):  ## Build snapshot release binaries and packages
	$(call title,Building snapshot artifacts)

	# create a config with the dist dir overridden
	echo "dist: $(SNAPSHOT_DIR)" > $(TEMP_DIR)/goreleaser.yml
	cat .goreleaser.yml >> $(TEMP_DIR)/goreleaser.yml

	# build release snapshots
	$(SNAPSHOT_CMD) --config $(TEMP_DIR)/goreleaser.yml

.PHONY: ci-release
ci-release: ci-check clean-dist
	$(call title,Publishing release artifacts)

	# create a config with the dist dir overridden
	echo "dist: $(DIST_DIR)" > $(TEMP_DIR)/goreleaser.yml
	cat .goreleaser.yml >> $(TEMP_DIR)/goreleaser.yml

	bash -c "\
		$(RELEASE_CMD) \
			--config $(TEMP_DIR)/goreleaser.yml \
				 || (cat /tmp/quill-*.log && false)"

.PHONY: ci-check
ci-check:
	ifndef $(CI)
		$(error "This step should ONLY be run in CI. Exiting...")
	endif

## Cleanup targets #################################

.PHONY: clean
clean: clean-dist clean-snapshot clean-images-docker
	$(call safe_rm_rf_children,$(TEMP_DIR))

.PHONY: clean-snapshot
clean-snapshot:
	$(call safe_rm_rf,$(SNAPSHOT_DIR))
	rm -f $(TEMP_DIR)/goreleaser.yml

.PHONY: clean-images-docker
clean-images-docker:
	docker images --format '{{.ID}} {{.Repository}}' | grep openshift4-mirror-go | awk '{print $$1}' | uniq | xargs -r docker rmi --force

.PHONY: clean-dist
clean-dist:
	$(call safe_rm_rf,$(DIST_DIR))
	rm -f $(TEMP_DIR)/goreleaser.yml
