#
# Copyright (c) 2019 Stefaan Coussement
# MIT License
#
# more info: https://github.com/stefaanc/terraform-provider-hosts
#
PLUGIN_NAME      := terraform-provider-hosts
PLUGIN_VERSION   := 0.0.0

OS               := $(shell go env GOOS)
ARCH             := $(shell go env GOARCH)
IS_WINDOWS       := $(findstring windows,$(OS))
RELEASE_PATH     := $(abspath ./releases/$(OS)_$(ARCH))
PLUGIN_PATH      := $(if $(IS_WINDOWS),$(APPDATA)/terraform.d/plugins,$(HOME)/.terraform.d/plugins)
PLUGIN_FILE      := $(PLUGIN_NAME)_v$(PLUGIN_VERSION)$(shell go env GOEXE)
TEST_LOG_FILE    := test.log
GO_TEST_FILES    := $(if $(IS_WINDOWS),$(shell dir /S /B *_test.go),$(shell find . -type f -name '*_test.go'))
GO_BUILD_FILES   := $(if $(IS_WINDOWS),$(shell dir /S /B *.go | findstr /v /c:"_test.go"),$(shell find . -type f -name '*.go' | grep -v '*_test.go'))

default: build

.PHONY: tidy
tidy:
	go mod tidy

$(TEST_LOG_FILE): $(GO_BUILD_FILES) $(GO_TEST_FILES)
ifneq (,$(IS_WINDOWS))
	PowerShell -NoProfile "go test ./... -v -cover | Tee-Object -FilePath \"$(TEST_LOG_FILE)\""
else
	go test ./... -v -cover | tee $(TEST_LOG_FILE)
endif

.PHONY: test
test: $(TEST_LOG_FILE)

$(PLUGIN_PATH)/$(PLUGIN_FILE): $(GO_BUILD_FILES)
	go build -o $(PLUGIN_PATH)/$(PLUGIN_FILE)

.PHONY: build
build: test $(PLUGIN_PATH)/$(PLUGIN_FILE)

$(RELEASE_PATH):
ifneq (,$(IS_WINDOWS))
	md $(subst /,\,$@)
else
	mkdir -p $@
endif

$(RELEASE_PATH)/$(PLUGIN_FILE): $(RELEASE_PATH) $(PLUGIN_PATH)/$(PLUGIN_FILE) $(TEST_LOG_FILE)
ifneq (,$(IS_WINDOWS))
	copy /Y $(subst /,\,$(PLUGIN_PATH))\$(PLUGIN_FILE) $(subst /,\,$(RELEASE_PATH))\$(PLUGIN_FILE)
	copy /Y $(TEST_LOG_FILE) $(subst /,\,$(RELEASE_PATH))\$(TEST_LOG_FILE)
else
	cp -f $(PLUGIN_PATH)/$(PLUGIN_FILE) $(RELEASE_PATH)/$(PLUGIN_FILE)
	cp -f $(TEST_LOG_FILE) $(RELEASE_PATH)/$(TEST_LOG_FILE)
endif

.PHONY: release
release: tidy test build $(RELEASE_PATH)/$(PLUGIN_FILE)
