.PHONY: release vet fmt test
release: vet fmt test

vet:
	go vet ./...

check:
	$(eval DIFFS:=$(shell goimports -l .))
	@if [ -n "$(DIFFS)" ]; then echo "goimports failed. Fix by running goimports."; echo "$(DIFFS)"; exit 1; fi

TESTDATA ?= $(shell pwd)/testdata

test:
	go test ./...
	@go run -tags=fizzbuzz $(TESTDATA) -fizzon 5 -buzzon 3	
	@go run -tags=example $(TESTDATA)
	@go run -tags=example $(TESTDATA) -c 18 -word "na " -newlines=false
	@go run -tags=example $(TESTDATA) -c 1 -word="Batman!"