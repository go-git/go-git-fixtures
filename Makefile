GOCMD = go
GOTEST = $(GOCMD) test 

test:
	$(GOTEST) ./...
