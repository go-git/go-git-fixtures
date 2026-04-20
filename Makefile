GOCMD = go
GOTEST = $(GOCMD) test 

# renovate: datasource=github-tags depName=golangci/golangci-lint
GOLANGCI_VERSION ?= v2.11.3
TOOLS_BIN := $(shell mkdir -p build/tools && realpath build/tools)

GOLANGCI = $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)
$(GOLANGCI):
	rm -f $(TOOLS_BIN)/golangci-lint*
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_VERSION)/install.sh | sh -s -- -b $(TOOLS_BIN) $(GOLANGCI_VERSION)
	mv $(TOOLS_BIN)/golangci-lint $(TOOLS_BIN)/golangci-lint-$(GOLANGCI_VERSION)

test:
	$(GOTEST) -race -parallel 20 ./...

validate: validate-lint validate-dirty validate-packs validate-bitmaps ## Run validation checks.

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
validate-packs:
	@find data -maxdepth 1 -type f -name 'pack-*.pack' | sort -u | \
	while read -r pack; do \
		base=$$(basename "$$pack" .pack); \
		hash=$${base#pack-}; \
		case "$${#hash}" in \
			40) \
				[ "$$hash" = "ee4fef0ef8be5053ebae4ce75acf062ddf3031fb" ] && continue; \
				git verify-pack -v "$$pack"; \
				git index-pack --rev-index -v "$$pack"; \
				;; \
			64) \
				[ "$$hash" = "407497645643e18a7ba56c6132603f167fe9c51c00361ee0c81d74a8f55d0ee2" ] && continue; \
				git --object-format=sha256 verify-pack -v "$$pack"; \
				git --object-format=sha256 index-pack --rev-index -v "$$pack"; \
				;; \
			*) \
				echo "Unknown hash length ($${#hash}) for $$pack" >&2; \
				exit 1; \
				;; \
		esac; \
	done
	@git status --short
	@git diff-index --quiet HEAD -- || { \
		echo "Generated pack metadata differs from HEAD" >&2; \
		exit 1; \
	}

validate-bitmaps:
	@find data -maxdepth 1 -type f -name '*.bitmap' | sort -u | \
	while read -r bitmap; do \
		pack=$${bitmap%.bitmap}.pack; \
		idx=$${bitmap%.bitmap}.idx; \
		if [ ! -f "$$pack" ]; then \
			echo "bitmap $$bitmap has no matching pack file" >&2; \
			exit 1; \
		fi; \
		if [ ! -f "$$idx" ]; then \
			echo "bitmap $$bitmap has no matching idx file" >&2; \
			exit 1; \
		fi; \
		base=$$(basename "$$pack" .pack); \
		hash=$${base#pack-}; \
		tmpdir=$$(mktemp -d); \
		case "$${#hash}" in \
			40) git init --bare "$$tmpdir" > /dev/null 2>&1 ;; \
			64) git init --bare --object-format=sha256 "$$tmpdir" > /dev/null 2>&1 ;; \
			*) echo "Unknown hash length ($${#hash}) for $$bitmap" >&2; rm -rf "$$tmpdir"; exit 1 ;; \
		esac; \
		cp "$$pack" "$$idx" "$$bitmap" "$$tmpdir/objects/pack/"; \
		commit=$$(GIT_DIR="$$tmpdir" git verify-pack -v "$$tmpdir/objects/pack/$$(basename $$idx)" 2>/dev/null \
			| awk '/commit/{print $$1; exit}'); \
		if [ -z "$$commit" ]; then \
			echo "no commit found in $$pack" >&2; rm -rf "$$tmpdir"; exit 1; \
		fi; \
		echo "$$commit" > "$$tmpdir/refs/heads/main"; \
		errout=$$(GIT_DIR="$$tmpdir" git rev-list --use-bitmap-index --count HEAD 2>&1 >/dev/null); \
		if [ -n "$$errout" ]; then \
			echo "$$bitmap: $$errout" >&2; rm -rf "$$tmpdir"; exit 1; \
		fi; \
		rm -rf "$$tmpdir"; \
		echo "OK: $$bitmap"; \
	done