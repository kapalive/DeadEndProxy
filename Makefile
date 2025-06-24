.PHONY: build deb clean release

APP := deadendproxy
VERSION := 1.0.5
ARCH := amd64
BUILD_DIR := build
DEB_DIR := $(BUILD_DIR)/deb
BIN_DIR := $(DEB_DIR)/usr/local/bin

build:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=0 go build -o $(BIN_DIR)/$(APP) ./cmd/main.go

deb: build
	mkdir -p $(DEB_DIR)/DEBIAN
	echo "Package: $(APP)" > $(DEB_DIR)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(DEB_DIR)/DEBIAN/control
	echo "Section: net" >> $(DEB_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	echo "Architecture: $(ARCH)" >> $(DEB_DIR)/DEBIAN/control
	echo "Maintainer: Devinsidercode" >> $(DEB_DIR)/DEBIAN/control
	echo "Description: Correct DeadEndProxy build with -config support" >> $(DEB_DIR)/DEBIAN/control

	# Создаём systemd unit
	mkdir -p $(DEB_DIR)/lib/systemd/system
	echo "[Unit]" > $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "Description=DeadEndProxy Service" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "After=network.target" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "[Service]" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "ExecStart=/usr/local/bin/deadendproxy -port-http 80 -port-proxy 443 -config /etc/deadendproxy/config.yaml" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "Restart=always" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "User=root" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "WorkingDirectory=/etc/deadendproxy" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "[Install]" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service
	echo "WantedBy=multi-user.target" >> $(DEB_DIR)/lib/systemd/system/deadendproxy.service

	dpkg-deb --build $(DEB_DIR) $(APP)_$(VERSION)_$(ARCH).deb

clean:
	rm -rf $(BUILD_DIR) *.deb

release: deb
	@echo "✅ DONE: deadendproxy_$(VERSION)_$(ARCH).deb"
