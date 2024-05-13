PLUGIN_SRC_DIR = ./internal/plugins/source/builtin
PLUGIN_COMPILED_DIR = ./internal/plugins/compiled
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

clean:
	rm -f $(PLUGIN_COMPILED_DIR)/*.so

plugins: clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildmode=plugin -o $(PLUGIN_COMPILED_DIR)/createplugin.so $(PLUGIN_SRC_DIR)/create-plugin/plugin.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildmode=plugin -o $(PLUGIN_COMPILED_DIR)/dateTime.so $(PLUGIN_SRC_DIR)/dateTime/plugin.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildmode=plugin -o $(PLUGIN_COMPILED_DIR)/memory.so $(PLUGIN_SRC_DIR)/memory/plugin.go

clara: plugins
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build 

run:
	go run ./cmd/main.go


