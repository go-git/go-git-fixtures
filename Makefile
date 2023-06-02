# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test 

# TODO: Move away from esc and into native Go embed.
esc:
	$(GOCMD) install github.com/mjibson/esc

test:
	$(GOTEST) ./...

generate: $(esc)
	$(GOCMD) generate
