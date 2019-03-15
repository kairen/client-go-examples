VERSION_MAJOR ?= 0
VERSION_MINOR ?= 3
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)
BUILD_IMAGE ?= gcr.io/google_containers/kube-cross:v1.9.2-1

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./out
$(shell mkdir -p $(BUILD_DIR))

ORG := github.com/kairen
REPOPATH ?= $(ORG)/simple-device-plugin
FILES := GOPATH=$(GOPATH) go list  -f '{{join .Deps "\n"}}' ./ | grep $(ORG) | GOPATH=$(GOPATH) xargs go list -f '{{ range $$file := .GoFiles }} {{$$.Dir}}/{{$$file}}{{"\n"}}{{end}}'

define DOCKER
	docker run --rm -e IN_DOCKER=1 --user $(shell id -u):$(shell id -g) -w /go/src/$(REPOPATH) -v $(GOPATH):/go --entrypoint /bin/bash $(1) -c '$(2)'
endef

ifeq ($(BUILD_IN_DOCKER),y)
	BUILD_IN_DOCKER=y
endif

ifneq ($(BUILD_OS),Linux)
	BUILD_IN_DOCKER=y
endif

ifeq ($(IN_DOCKER),1)
	BUILD_IN_DOCKER=n
endif

out/sdp: gopath out/sdp-$(GOOS)-$(GOARCH)
	mv $(BUILD_DIR)/sdp-$(GOOS)-$(GOARCH) $(BUILD_DIR)/sdp

out/sdp-%-amd64: $(shell $(FILES))
ifeq ($(BUILD_IN_DOCKER),y)
	$(call DOCKER,$(BUILD_IMAGE),/usr/bin/make $@)
else
	GOOS=$* go build -a -o $@ $(REPOPATH)
endif

.PHONY: all
all: cross

.PHONY: cross
cross: out/sdp-linux-amd64 out/sdp-darwin-amd64

.PHONY: deps
deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: build_image
build_image:
	docker build -t kairen/simple-device-plugin:$(VERSION) .

.PHONY: test
test:
	./test.sh

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: gopath
gopath:
ifneq ($(GOPATH)/src/$(REPOPATH),$(PWD))
	$(warning Warning: Building simple-device-plugin outside the GOPATH, should be $(GOPATH)/src/$(REPOPATH) but is $(PWD))
endif
