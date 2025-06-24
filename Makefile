.PHONY: build deb clean release

APP := deadendproxy
PKG := deadendproxy
ARCH := amd64
BUILD_DIR := build
SRC := ./cmd/main.go

RAW_VERSION := $(shell git describe --tags --always --dirty)
VERSION := $(shell echo $(RAW_VERSION) | sed 's/^v//')

BUILD_PATH := $(BUILD_DIR)/deb/usr/local/bin
DEB_DIR := $(BUILD_DIR)/deb
DEB_PACKAGE := $(PKG)_$(VERSION)_$(ARCH).deb

build:
	@echo "💻 Building $(APP) v$(VERSION)..."
	@mkdir -p $(BUILD_PATH)
	GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_PATH)/$(APP) $(SRC)

deb: build
	@echo "📦 Building .deb package: $(DEB_PACKAGE)"
	@mkdir -p $(DEB_DIR)/DEBIAN
	@printf "Package: $(PKG)\nVersion: $(VERSION)\nSection: net\nPriority: optional\nArchitecture: $(ARCH)\nMaintainer: Devinsidercode <dev@devinsidercode.com>\nDescription: DeadEndProxy - lightweight Go proxy server\n" > $(DEB_DIR)/DEBIAN/control
	@echo '#!/bin/sh' > $(DEB_DIR)/DEBIAN/postinst
	@echo 'set -e' >> $(DEB_DIR)/DEBIAN/postinst
	@echo 'systemctl restart deadendproxy || true' >> $(DEB_DIR)/DEBIAN/postinst
	@chmod 755 $(DEB_DIR)/DEBIAN/postinst
	@dpkg-deb --build $(DEB_DIR) $(DEB_PACKAGE)
	@echo "✅ Created package: $(DEB_PACKAGE)"

clean:
	@rm -rf $(BUILD_DIR)
	@rm -f $(DEB_PACKAGE)

release: deb
	@echo "🚀 Upload $(DEB_PACKAGE) to your apt repo directory"
