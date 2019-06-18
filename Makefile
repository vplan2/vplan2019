# BUILD TOOLS COMMANDS / LOCATION
GO      = go
DEP     = dep
GIT     = git
GOLINT  = golint
GREP    = grep
ZOLA    = zola
CLOC    = cloc
DOCKER  = docker

PACKAGE	= github.com/zekroTJA/vplan2019
GOPATH	= $(CURDIR)/.gopath
WDIR	= $(GOPATH)/src/$(PACKAGE)

BINNAME	= vplan2019_server
BINLOC	= $(CURDIR)

DTAG    = zekro/vplan2019:testing

ifeq ($(OS),Windows_NT)
	EXTENSION = .exe
endif

BIN		= $(BINLOC)/$(BINNAME)$(EXTENSION)

TAG		= $(shell $(GIT) describe --tags)
COMMIT	= $(shell $(GIT) rev-parse HEAD)
GOVERS  = $(shell $(GO) version | sed -e 's/ /_/g')

.PHONY: _make deps cleanup _finish run lint offline release \
	frontend frontend-vps cloc help crun release-vps

_make: $(WDIR) deps $(BIN) cleanup _finish

offline: $(WDIR) $(BIN) cleanup _finish

$(WDIR):
	@echo [ INFO ] creating working directory '$@'...
	mkdir -p $@
	cp -R $(CURDIR)/* $@/

$(BIN):
	@echo [ INFO ] building binary '$(BIN)'...
	(env GOPATH=$(GOPATH) $(GO) build -v -o $@ -ldflags "\
		-X $(PACKAGE)/internal/ldflags.AppVersion=$(TAG) \
		-X $(PACKAGE)/internal/ldflags.AppCommit=$(COMMIT) \
		-X $(PACKAGE)/internal/ldflags.GoVersion=$(GOVERS) \
		-X $(PACKAGE)/internal/ldflags.Release=TRUE" \
		$(WDIR)/cmd/server)

frontend:
	@echo [ INFO ] building frontend...
	cp $(CURDIR)/config/frontend.release.toml $(CURDIR)/web/config.toml
	cd $(CURDIR)/web && \
		$(ZOLA) build

frontend-vps:
	@echo [ INFO ] building frontend...
	cp $(CURDIR)/config/frontend.debug.toml $(CURDIR)/web/config.toml
	cd $(CURDIR)/web && \
		$(ZOLA) build

release: cleanup $(WDIR) deps frontend $(BIN) cleanup
	@echo [ INFO ] Creating release...
	mkdir $(CURDIR)/release
	mv -f $(BIN) $(CURDIR)/release
	[ "$(GOOS)" = "windows" ] && \
		mv $(CURDIR)/release/$(BINNAME) $(CURDIR)/release/$(BINNAME).exe || true
	cp -f -R $(CURDIR)/web/public $(CURDIR)/release/web

release-vps: cleanup $(WDIR) deps frontend-vps $(BIN) cleanup
	@echo [ INFO ] Creating release...
	mkdir $(CURDIR)/release
	mv -f $(BIN) $(CURDIR)/release
	[ "$(GOOS)" = "windows" ] && \
		mv $(CURDIR)/release/$(BINNAME) $(CURDIR)/release/$(BINNAME).exe || true
	cp -f -R $(CURDIR)/web/public $(CURDIR)/release/web

ferun: cleanup
	@echo [ INFO ] serving local frontend files...
	cp $(CURDIR)/config/frontend.local.toml $(CURDIR)/web/config.toml
	cd $(CURDIR)/web && \
		$(ZOLA) serve

deps:
	@echo [ INFO ] getting dependencies...	
	cd $(WDIR) && \
		(env GOPATH=$(GOPATH) $(DEP) ensure -v )

cleanup:
	@echo [ INFO ] cleaning up...
	[ -d $(GOPATH) ] && rm -r -f $(GOPATH) || true
	[ -d ./web/public ] && rm -r -f ./web/public || true
	

_finish:
	@echo ------------------------------------------------------------------------------
	@echo [ INFO ] Build successful.
	@echo [ INFO ] Your build is located at '$(BIN)'

run:
	@echo [ INFO ] Debug running...
	[ -d $(CURDIR)/web/public ] || { \
		cp $(CURDIR)/config/frontend.debug.toml $(CURDIR)/web/config.toml && \
		cd $(CURDIR)/web && \
		$(ZOLA) build; \
	}
	(env GOPATH=$(CURDIR)/../../../.. $(GO) run -v ./cmd/server -c $(CURDIR)/config/private.config.yml -web $(CURDIR)/web/public ${ARGS})

crun: cleanup run

lint:
	$(GOLINT) ./... | $(GREP) -v vendor

cloc:
	@echo [ INFO ] counting lines of code...
	$(CLOC) \
		--exclude-dir=vendor,docs,public \
		--exclude-lang=JSON,Markdown,YAML,XML,TOML,Sass ./

docker-build:
	docker build . -t $(DTAG)

docker-run:
	docker run \
            -p 8080:8080 \
            -v $(CURDIR)/testing/config:/etc/vplan/config \
            -v $(CURDIR)/testing/certs:/etc/vplan/certs \
            -d \
            $(DTAG)

help:
	@echo "Available recipes:"
	@echo "  *default* : created binary of backend server"
	@echo "  cleanup   : delete ./release, ./.gopath and ./web/public"
	@echo "  cloc      : count lines of code"
	@echo "  crun      : runs cleanup before run"
	@echo "  deps      : pulling dependencies from internet to ./vendor"
	@echo "  frontend  : compile frontend files to ./web/public"
	@echo "  lint      : go lint backend"
	@echo "  offline   : creates binaries witout pulling deps from the internet and using vendor instead"
	@echo "  release   : comple backend and frontend files to ./release"
	@echo "  run       : compile frontend file if not existent and go run backend"
