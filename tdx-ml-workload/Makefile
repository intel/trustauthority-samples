SHELL := /bin/bash

GITCOMMIT := $(shell git describe --always)
VERSION := v1.0.0
BUILDDATE := $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%S%z)

makefile_path := $(realpath $(lastword $(MAKEFILE_LIST)))
makefile_dir := $(dir $(makefile_path))
OUTDIR := $(addprefix $(makefile_dir),out)
TMPDIR := $(addprefix $(makefile_dir),tmp)

.PHONY: all installer clean help

all: clean installer

trustauthority-demo:
	CGO_CFLAGS_ALLOW="-f.*" CGO_CFLAGS="-O2 -D_FORTIFY_SOURCE=2" go build -buildmode=pie -trimpath -ldflags "-s -linkmode=external -extldflags -Wl,-O1,-z,relro,-z,now \
		-X github.com/intel/trustauthority-samples/tdxexample/version.BuildDate=$(BUILDDATE) -X github.com/intel/trustauthority-samples/tdxexample/version.Version=$(VERSION) \
		-X github.com/intel/trustauthority-samples/tdxexample/version.GitHash=$(GITCOMMIT)" -o trustauthority-demo

installer: trustauthority-demo
	mkdir -p installer
	cp build/linux/* installer/
	chmod +x installer/install.sh
	cp trustauthority-demo installer/trustauthority-demo
	makeself installer trustauthority-demo-$(VERSION).bin "trustauthority-demo $(VERSION)" ./install.sh
	rm -rf installer

clean:
	if pushd $(makefile_dir); then \
		rm -rf $(OUTDIR) $(TMPDIR); \
		rm -f *.bin trustauthority-demo; \
	fi;

help:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'