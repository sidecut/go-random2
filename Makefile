SWIFT_BUILD_FLAGS_DEBUG = -c debug
SWIFT_BUILD_FLAGS_RELEASE = -c release
BUILD_DIR = .build
BIN_NAME = go-random2
INSTALL_DIR = $(HOME)/.bin

.PHONY: all debug release install clean

all: debug

debug:
	swift build $(SWIFT_BUILD_FLAGS_DEBUG)

release:
	swift build $(SWIFT_BUILD_FLAGS_RELEASE)

install: release
	mkdir -p $(INSTALL_DIR)
	install -s $(BUILD_DIR)/release/$(BIN_NAME) $(INSTALL_DIR)/

clean:
	swift package clean
	rm -rf $(BUILD_DIR)
