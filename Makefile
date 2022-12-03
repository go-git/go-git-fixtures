GOCMD = go
GOTEST = $(GOCMD) test 

# Coverage
COVERAGE_REPORT = coverage.out
COVERAGE_MODE = count

test:
	@echo "running against `git version`"; \
	$(GOTEST) ./...

test-coverage:
	@echo "running against `git version`"; \
	echo "" > $(COVERAGE_REPORT); \
	$(GOTEST) -coverprofile=$(COVERAGE_REPORT) -coverpkg=./... -covermode=$(COVERAGE_MODE) ./...
