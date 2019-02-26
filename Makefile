BINARY = weed-exporter
package = github.com/andy/weed-exporter

SOURCE_DIR = github.com/andy/weed-exporter

all: build

.PHONY: all build deps clean

build: deps
		go build -o $(BINARY) $(SOURCE_DIR)

deps:
		go get -d $(package)

clean:
		go clean -i $(SOURCE_DIR)
			rm -rf $(BINARY)

