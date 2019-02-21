# Copyright © 2018 Kyle Bai(kyle.b@inwinstack.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)
BUILD_IMAGE ?= gcr.io/google_containers/kube-cross:v1.9.2-1

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./out
$(shell mkdir -p $(BUILD_DIR))

ORG := github.com/kairen
REPOPATH ?= $(ORG)/simple-operator

FILES := GOPATH=$(GOPATH) go list  -f '{{join .Deps "\n"}}' ./cmd/ | grep $(ORG) | GOPATH=$(GOPATH) xargs go list -f '{{ range $$file := .GoFiles }} {{$$.Dir}}/{{$$file}}{{"\n"}}{{end}}'

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

out/operator: gopath out/operator-$(GOOS)-$(GOARCH)
	mv $(BUILD_DIR)/operator-$(GOOS)-$(GOARCH) $(BUILD_DIR)/operator

out/operator-%-amd64: $(shell $(FILES))
ifeq ($(BUILD_IN_DOCKER),y)
	$(call DOCKER,$(BUILD_IMAGE),/usr/bin/make $@)
else
	GOOS=$* go build -ldflags="-X $(REPOPATH)/pkg/version.version=$(VERSION)" -a -o $@ $(REPOPATH)/cmd
endif

.PHONY: all
all: cross

.PHONY: cross
cross: out/operator-linux-amd64 out/operator-darwin-amd64

.PHONY: codegen
codegen:
	./hack/codegen.sh

.PHONY: deps
deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: build_image
build_image: deps
	docker build -t kairen/simple-operator:$(VERSION) .

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: gopath
gopath:
ifneq ($(GOPATH)/src/$(REPOPATH),$(PWD))
	$(warning Warning: Building operator outside the GOPATH, should be $(GOPATH)/src/$(REPOPATH) but is $(PWD))
endif
