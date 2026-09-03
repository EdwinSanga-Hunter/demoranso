.PHONY: all deps

all: build clean

PROJECT_DIR=$(shell pwd)
BUILD_DIR=$(PROJECT_DIR)/build
BIN_DIR=$(PROJECT_DIR)/bin
TOOLS_BIN=$(shell go env GOPATH)/bin

SERVER_HOST?=localhost
SERVER_PORT?=8080
SERVER_URL=http://$(SERVER_HOST):$(SERVER_PORT)

# Target OS/arch for the ransomware and unlocker binaries.
# Defaults to the host platform so you can run the demo locally.
# For the classic windows .exe files: make -e CLIENT_OS=windows CLIENT_ARCH=386
CLIENT_OS?=$(shell go env GOOS)
ifeq ($(CLIENT_OS),windows)
CLIENT_ARCH?=386
CLIENT_EXT=.exe
else
CLIENT_ARCH?=$(shell go env GOARCH)
endif

# Target OS for the server binary (GOOS kept for backward compatibility:
# make -e GOOS=windows builds the server for windows)
SERVER_OS?=$(GOOS)
ifeq ($(SERVER_OS),)
SERVER_OS=$(shell go env GOOS)
endif
ifeq ($(SERVER_OS),windows)
SERVER_EXT=.exe
endif

LINKER_VARS=-X main.ServerBaseURL=$(SERVER_URL) -X main.UseTor=$(USE_TOR)
ifeq ($(CLIENT_OS),windows)
CLIENT_LDFLAGS=-s -w $(HIDDEN) $(LINKER_VARS)
else
CLIENT_LDFLAGS=-s -w $(LINKER_VARS)
endif

deps:
	go mod download
	GOBIN=$(TOOLS_BIN) go install github.com/akavel/rsrc@latest
	GOBIN=$(TOOLS_BIN) go install github.com/go-bindata/go-bindata/v3/go-bindata@latest

pre-build: clean-bin
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/ransomware $(BUILD_DIR)/unlocker $(BUILD_DIR)/server $(BIN_DIR)/server
	openssl genrsa -traditional -out $(BUILD_DIR)/server/private.pem 4096
	openssl rsa -in $(BUILD_DIR)/server/private.pem -outform PEM -pubout -out $(BUILD_DIR)/ransomware/public.pem
	cd $(BUILD_DIR)/ransomware && $(TOOLS_BIN)/go-bindata -pkg main -o public_key.go public.pem
	cp -r cmd/ransomware cmd/unlocker cmd/server $(BUILD_DIR)/
	openssl req -x509 -newkey rsa:2048 -keyout $(BUILD_DIR)/server/key.pem -out $(BUILD_DIR)/server/cert.pem -days 365 -nodes -subj "/CN=$(SERVER_HOST)" >/dev/null 2>&1
ifeq ($(CLIENT_OS),windows)
	$(TOOLS_BIN)/rsrc -arch=$(CLIENT_ARCH) -manifest ransomware.manifest -ico icon.ico -o $(BUILD_DIR)/ransomware/ransomware.syso
	cp $(BUILD_DIR)/ransomware/ransomware.syso $(BUILD_DIR)/unlocker/unlocker.syso
endif

build: pre-build
	cd $(BUILD_DIR)/ransomware && GOOS=$(CLIENT_OS) GOARCH=$(CLIENT_ARCH) go build -ldflags "$(CLIENT_LDFLAGS)" -o $(BIN_DIR)/ransomware$(CLIENT_EXT)
	cd $(BUILD_DIR)/unlocker && GOOS=$(CLIENT_OS) GOARCH=$(CLIENT_ARCH) go build -ldflags "-s -w" -o $(BIN_DIR)/unlocker$(CLIENT_EXT)
	cd $(BUILD_DIR)/server && GOOS=$(SERVER_OS) go build -ldflags "-s -w" -o $(BIN_DIR)/server/server$(SERVER_EXT)
	cp $(BUILD_DIR)/server/private.pem $(BUILD_DIR)/server/cert.pem $(BUILD_DIR)/server/key.pem $(BIN_DIR)/server/

# Create some dummy files in the default demo folder (~/ransomware-demo)
demo-files:
	mkdir -p ~/ransomware-demo
	printf 'Top secret demo content\n' > ~/ransomware-demo/secret.txt
	printf '%%PDF-1.4 demo' > ~/ransomware-demo/report.pdf
	printf 'PNGDATA' > ~/ransomware-demo/photo.png
	printf 'DBDATA' > ~/ransomware-demo/database.db
	@echo "Demo files created in ~/ransomware-demo"

# Create double-click launchers on the Desktop (linux demo, no terminal needed)
demo-launcher:
	@mkdir -p ~/Desktop
	@printf '%s\n' '[Desktop Entry]' 'Type=Application' 'Name=Demo Server' 'Comment=Start the demo server' 'Exec=x-terminal-emulator -e $(PROJECT_DIR)/scripts/run-server.sh' 'Terminal=true' 'Icon=applications-system' > ~/Desktop/demo-server.desktop
	@printf '%s\n' '[Desktop Entry]' 'Type=Application' 'Name=Ransomware Demo' 'Comment=Educational ransomware demo' 'Exec=x-terminal-emulator -e $(PROJECT_DIR)/bin/ransomware' 'Terminal=true' 'Icon=applications-system' > ~/Desktop/ransomware-demo.desktop
	@printf '%s\n' '[Desktop Entry]' 'Type=Application' 'Name=Unlocker Demo' 'Comment=Educational unlocker demo' 'Exec=x-terminal-emulator -e $(PROJECT_DIR)/bin/unlocker' 'Terminal=true' 'Icon=applications-system' > ~/Desktop/unlocker-demo.desktop
	@chmod +x ~/Desktop/demo-server.desktop ~/Desktop/ransomware-demo.desktop ~/Desktop/unlocker-demo.desktop
	@echo "Launchers created on ~/Desktop (first run: right-click -> Allow Launching)"

clean:
	rm -r $(BUILD_DIR) || true

clean-bin:
	rm -r $(BIN_DIR) || true
