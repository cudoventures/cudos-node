BUILDDIR ?= $(CURDIR)/build

all: install

$(BUILDDIR)/:
	mkdir -p $@

BUILD_TARGETS := build install
build: $(BUILDDIR)/
build: BUILD_ARGS=-o=$(BUILDDIR)

$(BUILD_TARGETS):
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

test:
	go test ./...

clean:
	rm -rf $(BUILDDIR)/

.PHONY: all $(BUILD_TARGETS) clean
