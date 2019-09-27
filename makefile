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
TEST_LOG_FILE    := _test.log
COVER_TMP_FILE   := _coverage.tmp
COVER_LOG_FILE   := _coverage.log
COVER_HTML_FILE  := _coverage.html
GO_TEST_FILES    := $(if $(IS_WINDOWS),$(shell dir /S /B *_test.go),$(shell find . -type f -name '*_test.go'))
GO_BUILD_FILES   := $(if $(IS_WINDOWS),$(shell dir /S /B *.go | findstr /v /c:"_test.go"),$(shell find . -type f -name '*.go' | grep -v '*_test.go'))

.PHONY: tidy       # tidy the module definition
.PHONY: test       # test the module and generate a test log and coverage report
.PHONY: log        # write the test log for the module to stdout
.PHONY: report     # write the coverage report for the module to stdout
.PHONY: browse     # open browser to analyse coverage for the module (only on windows)
.PHONY: build      # build the module
.PHONY: release    # release the module

.PHONY: default    # change to the actions you want to run by default
default: test log

tidy:
	go mod tidy

#
# testing

$(COVER_TMP_FILE) $(TEST_LOG_FILE): $(GO_BUILD_FILES) $(GO_TEST_FILES)
ifneq (,$(IS_WINDOWS))
	PowerShell -NoProfile "go test ./... -v -coverprofile $(COVER_TMP_FILE) > $(TEST_LOG_FILE); if ( $$LASTEXITCODE -ne 0 ) { Get-Content $(TEST_LOG_FILE) }; exit $$LASTEXITCODE"
else
	go test ./... -v -coverprofile $(COVER_TMP_FILE) > $(TEST_LOG_FILE); if [ $? != 0 ] ; then cat $(TEST_LOG_FILE) ; fi
endif

$(COVER_LOG_FILE) $(COVER_HTML_FILE): $(COVER_TMP_FILE)
ifneq (,$(IS_WINDOWS))
	PowerShell -NoProfile "go tool cover -func $(COVER_TMP_FILE) -o $(COVER_LOG_FILE)"
	PowerShell -NoProfile "go tool cover -html $(COVER_TMP_FILE) -o $(COVER_HTML_FILE)"
#	del /f $(COVER_TMP_FILE)
else
	go tool cover -func $(COVER_TMP_FILE) -o $(COVER_LOG_FILE)
	go tool cover -html $(COVER_TMP_FILE) -o $(COVER_HTML_FILE)
#	rm -f $(COVER_TMP_FILE)
endif

test: $(TEST_LOG_FILE)

log: $(TEST_LOG_FILE)
ifneq (,$(IS_WINDOWS))
	PowerShell -NoProfile "Get-Content $(TEST_LOG_FILE)"
else
	cat $(TEST_LOG_FILE)
endif

report: $(COVER_LOG_FILE)
ifneq (,$(IS_WINDOWS))
	PowerShell -NoProfile "Get-Content $(COVER_LOG_FILE)"
else
	cat $(COVER_LOG_FILE)
endif

browse: $(COVER_HTML_FILE)
ifneq (,$(IS_WINDOWS))
	timout /t 5; start $(COVER_HTML_FILE)
endif

#
# building

$(PLUGIN_PATH)/$(PLUGIN_FILE): $(GO_BUILD_FILES)
	go build -o $(PLUGIN_PATH)/$(PLUGIN_FILE)

build: $(PLUGIN_PATH)/$(PLUGIN_FILE)

#
# releasing

$(RELEASE_PATH):
ifneq (,$(IS_WINDOWS))
	md $(subst /,\,$@)
else
	mkdir -p $@
endif

$(RELEASE_PATH)/$(PLUGIN_FILE): $(RELEASE_PATH) $(TEST_LOG_FILE) $(COVER_LOG_FILE) $(COVER_HTML_FILE) $(PLUGIN_PATH)/$(PLUGIN_FILE)
ifneq (,$(IS_WINDOWS))
	copy /Y $(TEST_LOG_FILE) $(subst /,\,$(RELEASE_PATH))\\$(PLUGIN_NAME)_v$(PLUGIN_VERSION)_test.log
	copy /Y $(COVER_LOG_FILE) $(subst /,\,$(RELEASE_PATH))\\$(PLUGIN_NAME)_v$(PLUGIN_VERSION)_coverage.log
	copy /Y $(COVER_HTML_FILE) $(subst /,\,$(RELEASE_PATH))\\$(PLUGIN_NAME)_v$(PLUGIN_VERSION)_coverage.html
	copy /Y $(subst /,\,$(PLUGIN_PATH))\$(PLUGIN_FILE) $(subst /,\,$(RELEASE_PATH))\\$(PLUGIN_FILE)
else
	cp -f $(TEST_LOG_FILE) $(RELEASE_PATH)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)_test.log
	cp -f $(COVER_LOG_FILE) $(RELEASE_PATH)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)_coverage.log
	cp -f $(COVER_HTML_FILE) $(RELEASE_PATH)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)_coverage.html
	cp -f $(PLUGIN_PATH)/$(PLUGIN_FILE) $(RELEASE_PATH)/$(PLUGIN_FILE)
endif

release: tidy test report build $(RELEASE_PATH)/$(PLUGIN_FILE)
