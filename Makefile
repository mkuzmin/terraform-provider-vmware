TEST?=$$(go list ./... | grep -v '/vendor/')
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test vet

# bin generates the releaseable binaries for Terraform
bin: fmtcheck
	@TF_RELEASE=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing Terraform locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: fmtcheck
	@TF_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

quickdev:
	@TF_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# test runs the unit tests
test: fmtcheck errcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4


# testrace runs the race checker
testrace: fmtcheck
	TF_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v /vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: bin core-dev core-test cover default dev errcheck fmt fmtcheck plugin-dev quickdev test testrace tools vet
