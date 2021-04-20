.PHONY: all
all: build

NOW 		:=	$(shell date -u +'%Y-%m-%dT%TZ')
REVISION 	:=	$(shell git rev-parse --short=8 HEAD)

LDFLAGS		:=	"-X main.revision=$(REVISION) -X main.buildTime=$(NOW)"

MAIN_CLI	:=	bin/sql-executor

.PHONY: build
build:
	go build -o $(MAIN_CLI) -ldflags $(LDFLAGS) cmd/executor/main.go

.PHONY: modclean
modclean:
	go mod tidy

.PHONY: clean
clean:
	go clean -testcache
	rm -f $(MAIN_CLI)

.PHONY: fmt
fmt:
	gofmt -w **/*.go

.PHONY: test
test:
	go test -v -race -cover ./...

GOLANGCI_LINT_CLI		:= bin/golangci-lint
GOLANGCI_LINT_SITE		:= https://raw.githubusercontent.com
GOLANGCI_LINT_PATH		:= /golangci/golangci-lint/master/install.sh
GOLANGCI_LINT_VERSION 	:= v1.39.0

$(GOLANGCI_LINT_CLI):
	@echo "# install $@ before lint check"
	@mkdir -p bin
	$(shell curl -sSfL $(GOLANGCI_LINT_SITE)$(GOLANGCI_LINT_PATH) | sh -s $(GOLANGCI_LINT_VERSION))

.PHONY: lint
lint: $(GOLANGCI_LINT_CLI)
	./$(GOLANGCI_LINT_CLI) run $(linters)

GOCREDITS	:= bin/gocredits

$(GOCREDITS):
	@echo "# install $@ before creating CREDITS file"
	@mkdir -p bin
	GOBIN=$(abspath bin) go get -u github.com/Songmu/gocredits/cmd/gocredits

.PHONY: CREDITS
CREDITS: $(GOCREDITS)
	@echo "# create CREDITS file"
	rm -f $@
	$(GOCREDITS) -skip-missing . > $@
