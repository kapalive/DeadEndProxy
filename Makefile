.PHONY: build deb clean release

APP := deadendproxy
VERSION := 1.0.5
ARCH := amd64
BUILD_DIR := build
DEB_DIR := $(BUILD_DIR)/deb
BIN_DIR := $(DEB_DIR)/usr/local/bin
CONFIG_DIR := $(DEB_DIR)/etc/deadendproxy
SYSTEMD_DIR := $(DEB_DIR)/lib/systemd/system

build:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=0 go build -o $(BIN_DIR)/$(APP) ./cmd/main.go

deb: build
	# ==== CONTROL FILE ====
	mkdir -p $(DEB_DIR)/DEBIAN
	echo "Package: $(APP)" > $(DEB_DIR)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(DEB_DIR)/DEBIAN/control
	echo "Section: net" >> $(DEB_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	echo "Architecture: $(ARCH)" >> $(DEB_DIR)/DEBIAN/control
	echo "Maintainer: Devinsidercode" >> $(DEB_DIR)/DEBIAN/control
	echo "Description: DeadEndProxy with working systemd and config.yaml" >> $(DEB_DIR)/DEBIAN/control

	# ==== SYSTEMD SERVICE ====
	mkdir -p $(SYSTEMD_DIR)
	echo "[Unit]" > $(SYSTEMD_DIR)/deadendproxy.service
	echo "Description=DeadEndProxy Service" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "After=network.target" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "[Service]" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "ExecStart=/usr/local/bin/deadendproxy -port-http 80 -port-proxy 443 -config /etc/deadendproxy/config.yaml" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "Restart=always" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "User=root" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "WorkingDirectory=/etc/deadendproxy" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "[Install]" >> $(SYSTEMD_DIR)/deadendproxy.service
	echo "WantedBy=multi-user.target" >> $(SYSTEMD_DIR)/deadendproxy.service

	# ==== CONFIG.YAML ====
	mkdir -p $(CONFIG_DIR)
	cp ./config.yaml $(CONFIG_DIR)/config.yaml

	# ==== DEB PACKAGE ====
	dpkg-deb --build $(DEB_DIR) $(APP)_$(VERSION)_$(ARCH).deb

clean:
	rm -rf $(BUILD_DIR) *.deb

release: deb
	@echo "✅ DONE: $(APP)_$(VERSION)_$(ARCH).deb"
