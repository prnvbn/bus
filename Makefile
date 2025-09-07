.PHONY : dist windows darwin linux-arm64 linux-amd64 clean

DIST_DIR="dist"
ZIP="zip -m"

LDFLAGS="-s -w"
BUILD_FLAGS=-ldflags=$(LDFLAGS)

dist:
	mkdir -p $(DIST_DIR)
	$(MAKE) windows darwin-arm64 darwin-amd64 linux-arm64 linux-amd64

.PHONY= build-windows
windows:
	CGO_ENABLED=0 GOOS=windows \
	go build $(BUILD_FLAGS) -o $(DIST_DIR)/bus-windows.exe

darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
	go build $(BUILD_FLAGS) -o $(DIST_DIR)/bus-darwin-amd64
darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
	go build $(BUILD_FLAGS) -o $(DIST_DIR)/bus-darwin-arm64

linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
	go build $(BUILD_FLAGS) -o  $(DIST_DIR)/bus-linux-arm64

linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build $(BUILD_FLAGS) -o  $(DIST_DIR)/bus-linux-amd64

clean:
	rm -r $(DIST_DIR)/