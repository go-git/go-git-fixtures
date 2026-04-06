GOCMD = go
GOTEST = $(GOCMD) test 

# renovate: datasource=github-tags depName=golangci/golangci-lint
GOLANGCI_VERSION ?= v2.10.1
TOOLS_BIN := $(shell mkdir -p build/tools && realpath build/tools)

GOLANGCI = $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)
$(GOLANGCI):
	rm -f $(TOOLS_BIN)/golangci-lint*
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_VERSION)/install.sh | sh -s -- -b $(TOOLS_BIN) $(GOLANGCI_VERSION)
	mv $(TOOLS_BIN)/golangci-lint $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)

test:
	$(GOTEST) -race -parallel 20 ./...

validate: validate-lint validate-dirty validate-packs ## Run validation checks.

validate-lint: $(GOLANGCI)
	$(GOLANGCI) run

validate-dirty:
ifneq ($(shell git status --porcelain --untracked-files=no),)
	@echo worktree is dirty
	@git --no-pager status
	@git --no-pager diff
	@exit 1
endif

validate-packs:
	@find data -maxdepth 1 -type f -name 'pack-*.pack' | sort | sort -u | \
	while read -r pack; do \
		base=$$(basename "$$pack" .pack); \
		hash=$${base#pack-}; \
		case "$${#hash}" in \
			40) \
 				[ "$$hash" = "ee4fef0ef8be5053ebae4ce75acf062ddf3031fb" ] && continue; \
				git verify-pack -v "$$pack"; \
				;; \
			64) \
				git verify-pack --object-format=sha256 -v "$$pack"; \
				;; \
			*) \
				echo "Unknown hash length ($${#hash}) for $$pack" >&2; \
				exit 1; \
				;; \
		esac; \
	done
