PLUGIN_NAME="caddy-trace"
PLUGIN_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
LATEST_GIT_COMMIT:=$(shell git log --format="%H" -n 1 | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
BUILD_DIR:=$(shell pwd)
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif
CADDY_VERSION="v2.6.4"

all:
	@echo "Version: $(PLUGIN_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"
	@mkdir -p bin/
	@rm -rf ./bin/caddy
	@rm -rf ../xcaddy-$(PLUGIN_NAME)/*
	@mkdir -p ../xcaddy-$(PLUGIN_NAME) && cd ../xcaddy-$(PLUGIN_NAME) && \
		xcaddy build $(CADDY_VERSION) --output ../$(PLUGIN_NAME)/bin/caddy \
		--with github.com/greenpau/caddy-trace@$(LATEST_GIT_COMMIT)=$(BUILD_DIR)
	@#bin/caddy run -environ -config assets/conf/config.json

.PHONY: linter
linter:
	@echo "Running lint checks"
	@golint *.go
	@echo "$@: complete"

.PHONY: test
test: covdir linter
	@echo "$@: starting"
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./*.go
	@echo "$@: complete"

.PHONY: ctest
ctest: covdir linter
	@echo "$@: starting"
	@time richgo test $(VERBOSE) $(TEST) -coverprofile=.coverage/coverage.out ./*.go
	@echo "$@: complete"

.PHONY: covdir
covdir:
	@echo "$@: starting"
	@echo "Creating .coverage/ directory"
	@mkdir -p .coverage
	@echo "$@: complete"

.PHONY: coverage
coverage:
	@echo "$@: starting"
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@go test -covermode=count -coverprofile=.coverage/coverage.out ./*.go
	@go tool cover -func=.coverage/coverage.out | grep -v "100.0"
	@echo "$@: complete"

.PHONY: docs
docs:
	@echo "$@: starting"
	@mkdir -p .doc
	@go doc -all > .doc/index.txt
	@echo "$@: complete"

.PHONY: clean
clean:
	@echo "$@: starting"
	@rm -rf .doc
	@rm -rf .coverage
	@rm -rf bin/
	@echo "$@: complete"

.PHONY: qtest
qtest: covdir
	@echo "$@: starting"
	@echo "Perform quick tests ..."
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out -run TestCaddyfile ./*.go
	@echo "$@: complete"

.PHONY: dep
dep:
	@echo "$@: starting"
	@echo "Making dependencies check ..."
	@go install golang.org/x/lint/golint@latest
	@go install github.com/kyoh86/richgo@latest
	@go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
	@go install github.com/greenpau/versioned/cmd/versioned@latest
	@echo "$@: complete"

.PHONY: release
release:
	@echo "$@: starting"
	@echo "Making release"
	@go mod tidy
	@go mod verify
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && false )
	@versioned -patch
	@echo "Patched version"
	@git add VERSION
	@git commit -m "released v`cat VERSION | head -1`"
	@git tag -a v`cat VERSION | head -1` -m "v`cat VERSION | head -1`"
	@git push
	@git push --tags
	@@echo "If necessary, run the following commands:"
	@echo "  git push --delete origin v$(PLUGIN_VERSION)"
	@echo "  git tag --delete v$(PLUGIN_VERSION)"
	@echo "$@: complete"

.PHONY: logo
logo:
	@echo "$@: starting"
	@mkdir -p assets/docs/images
	@gm convert -background black -font Bookman-Demi \
		-size 640x320 "xc:black" \
		-pointsize 72 \
		-draw "fill white gravity center text 0,0 'caddy\ntrace'" \
		assets/docs/images/logo.png
	@echo "$@: complete"
