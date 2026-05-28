SWIFT_BUILD_FLAGS_DEBUG = -c debug
SWIFT_BUILD_FLAGS_RELEASE = -c release
BUILD_DIR = .build
BIN_NAME = rndutil
INSTALL_DIR = $(HOME)/.bin
GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
VERSION_FILE = Sources/Version.swift

.PHONY: all debug release install clean version

all: debug

$(VERSION_FILE):
	printf 'enum Build {\n    static let version = "$(GIT_VERSION)"\n}\n' > $(VERSION_FILE)

version:
	printf 'enum Build {\n    static let version = "$(GIT_VERSION)"\n}\n' > $(VERSION_FILE)

debug: version
	swift build $(SWIFT_BUILD_FLAGS_DEBUG)

release: version
	swift build $(SWIFT_BUILD_FLAGS_RELEASE)

install: release
	mkdir -p $(INSTALL_DIR)
	install -s $(BUILD_DIR)/release/$(BIN_NAME) $(INSTALL_DIR)/

clean:
	swift package clean
	rm -rf $(BUILD_DIR)
	rm -f $(VERSION_FILE)
