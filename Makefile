# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test 

test:
	$(GOTEST) ./...

generate: $(esc)
	$(GOCMD) generate
