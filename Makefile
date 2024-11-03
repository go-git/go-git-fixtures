GOCMD = go
GOTEST = $(GOCMD) test 

GOLANGCI_VERSION ?= v1.61.0
TOOLS_BIN := $(shell mkdir -p build/tools && realpath build/tools)

GOLANGCI = $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)
$(GOLANGCI):
	rm -f $(TOOLS_BIN)/golangci-lint*
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_VERSION)/install.sh | sh -s -- -b $(TOOLS_BIN) $(GOLANGCI_VERSION)
	mv $(TOOLS_BIN)/golangci-lint $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)

test:
	$(GOTEST) -race -parallel 20 ./...

validate: validate-lint validate-dirty ## Run validation checks.

validate-lint: $(GOLANGCI)
	$(GOLANGCI) run

define go-install-tool
@[ -f $(1) ] || { \
set -e ;\
echo "Downloading $(2)" ;\
GOBIN=$(TOOLS_BIN) go install $(2) ;\
}
endef

validate-dirty:
ifneq ($(shell git status --porcelain --untracked-files=no),)
	@echo worktree is dirty
	@git --no-pager status
	@git --no-pager diff
	@exit 1
endif
