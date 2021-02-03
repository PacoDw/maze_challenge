GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TEST?=./...
GOLANGCI_VERSION=v1.35.2
MISSPELL_VERSION=v0.3.4

# test all the existing test files or just one
test: fmtcheck
	@if [ "${name}" = "" ] ; then \
		go test $(TEST) -v -timeout=30s -parallel=4 ; \
	else \
	    go test $(TEST) -v -run $(name) -timeout=30s -parallel=4 ; \
	fi

clean-cache:
	@go clean -cache -modcache -i -r

# check what files need to format
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

# check if the tools are installed
checktools:
	@sh -c "'$(CURDIR)/scripts/checktools.sh'"

# format all the files with gofmt
fmt:
	@if [ "${pkg}" = "" ] ; then \
		echo "Error: please add a package to be formatted, e.g.:" ; \
		echo "make fmt pkg=kafka" ; \
	else \
		echo "==> Fixing source code with gofmt..." ; \
		gofmt -s -w ./$(pkg) ; \
	fi

# Install dev lints tools 
tools:
	@echo ""
	@echo "==> Installing missing commands dependencies..."
	curl -sSfL https://raw.githubusercontent.com/client9/misspell/master/install-misspell.sh | sh -s $(MISSPELL_VERSION)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_VERSION)

# it runs the lints
lint: fmt
## this is the same that make checktools command to validate the tools
	@if ! sh -c "'$(CURDIR)/scripts/checktools.sh'" 2>&1 /dev/null ; then \
		$(MAKE) -s tools ; \
	fi
	@if [ "${pkg}" = "" ] ; then \
		echo "Error: please add a package to be tested wiht the lint, e.g.:" ; \
		echo "make lint pkg=kafka" ; \
	else \
		echo "" ; \
		$(MAKE) fmt $(pkg) ; \
		echo "==> Checking source code against linters..." ; \
		bin/golangci-lint run ./$(pkg) -v ; \
	fi


# check all the test and run the linter
check: test lint

.PHONY: test fmtcheck checktools fmt tools lint check