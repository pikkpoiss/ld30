.phony: build clean run

PROJECT = ld30
SOURCES = $(wildcard src/*.go)
RUNTIME_ASSETS = $(wildcard src/assets/*)
ICON_ASSETS = $(wildcard assets/*.icns)

BASEBUILD = build/$(PROJECT)-osx

OSXBUILD = $(BASEBUILD)/$(PROJECT).app/Contents
OSXLIBS  = $(wildcard libs/osx/*.dylib)
OSXLIBSD = $(subst libs/osx/,$(OSXBUILD)/MacOS/,$(OSXLIBS))

WINBUILD = build/$(PROJECT)-win
WINLIBS  = $(wildcard libs/win/*.dll)
WINLIBSD = $(subst libs/win/,$(WINBUILD)/,$(WINLIBS))

VERSION = $(shell cat VERSION)
REPLACE = s/9\.9\.9/$(VERSION)/g

clean:
	rm -rf build

$(OSXBUILD)/MacOS/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/Info.plist: pkg/osx/Info.plist
	mkdir -p $(OSXBUILD)
	sed $(REPLACE) $< > $@

$(OSXBUILD)/MacOS/%.dylib: libs/osx/%.dylib
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go
#	cd $(OSXBUILD)/MacOS/ && ../../../../../scripts/fix.sh

$(OSXBUILD)/Resources/%.icns: assets/%.icns
	mkdir -p $(dir $@)
	cp $< $@
$(OSXBUILD)/Resources/assets/%: src/assets/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-osx-$(VERSION).zip: \
	$(OSXBUILD)/MacOS/launch.sh \
	$(OSXBUILD)/Info.plist \
	$(OSXLIBSD) \
	$(OSXBUILD)/MacOS/$(PROJECT) \
	$(subst src/assets/,$(OSXBUILD)/Resources/assets/,$(RUNTIME_ASSETS)) \
	$(subst assets/,$(OSXBUILD)/Resources/,$(ICON_ASSETS))
	cd build && zip -r $(notdir $@) $(PROJECT)-osx

$(WINBUILD)/$(PROJECT).exe: $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

$(WINBUILD)/%.dll: libs/win/%.dll
	mkdir -p $(dir $@)
	cp $< $@

$(WINBUILD)/assets/%: src/assets/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-win-$(VERSION).zip: \
	$(WINBUILD)/$(PROJECT).exe \
	$(WINLIBSD) \
	$(subst src/assets/,$(WINBUILD)/Resources/assets/,$(RUNTIME_ASSETS)) \
	cd build && /c/Program\ Files/7-Zip/7z.exe a -r $(notdir $@) $(PROJECT)-win

build: build/$(PROJECT)-osx-$(VERSION).zip

run: build
	$(OSXBUILD)/MacOS/launch.sh
