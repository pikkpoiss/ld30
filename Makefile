.phony: build clean run

PROJECT = sol
SOURCES = $(wildcard src/*.go)
RUNTIME_ASSETS = $(wildcard src/assets/*)
ICON_ASSETS = $(wildcard assets/*.icns)

BASEBUILD = build/$(PROJECT)-osx

OSXBUILD = $(BASEBUILD)/$(PROJECT).app/Contents
OSXLIBS  = $(wildcard libs/osx/*.dylib)
OSXLIBSD = $(subst libs/osx/,$(OSXBUILD)/MacOS/,$(OSXLIBS))

YOSBUILD = $(BASEBUILD)-yosemite/$(PROJECT).app/Contents
YOSLIBS  = $(wildcard libs/osx-yosemite/*.dylib)
YOSLIBSD = $(subst libs/osx-yosemite/,$(YOSBUILD)/MacOS/,$(YOSLIBS))

WINBUILD = build/$(PROJECT)-win
WINLIBS  = $(wildcard libs/win/*.dll)
WINLIBSD = $(subst libs/win/,$(WINBUILD)/,$(WINLIBS))

NIXBUILD = build/$(PROJECT)-linux
NIXLIBS  = $(wildcard libs/linux/*.*)
NIXLIBSD = $(subst libs/linux/,$(NIXBUILD)/libs/,$(NIXLIBS))

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
	cd $(OSXBUILD)/MacOS/ && ../../../../../scripts/fix.sh

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

$(YOSBUILD)/MacOS/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(YOSBUILD)/Info.plist: pkg/osx/Info.plist
	mkdir -p $(YOSBUILD)
	sed $(REPLACE) $< > $@

$(YOSBUILD)/MacOS/%.dylib: libs/osx-yosemite/%.dylib
	mkdir -p $(dir $@)
	cp $< $@

$(YOSBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go
	cd $(YOSBUILD)/MacOS/ && ../../../../../scripts/fix-yosemite.sh

$(YOSBUILD)/Resources/%.icns: assets/%.icns
	mkdir -p $(dir $@)
	cp $< $@

$(YOSBUILD)/Resources/assets/%: src/assets/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-osx-yosemite-$(VERSION).zip: \
	$(YOSBUILD)/MacOS/launch.sh \
	$(YOSBUILD)/Info.plist \
	$(YOSLIBSD) \
	$(YOSBUILD)/MacOS/$(PROJECT) \
	$(subst src/assets/,$(YOSBUILD)/Resources/assets/,$(RUNTIME_ASSETS)) \
	$(subst assets/,$(YOSBUILD)/Resources/,$(ICON_ASSETS))
	cd build && zip -r $(notdir $@) $(PROJECT)-osx-yosemite

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
	$(subst src/assets/,$(WINBUILD)/assets/,$(RUNTIME_ASSETS))
	cd build && /c/Program\ Files/7-Zip/7z.exe a -r $(notdir $@) $(PROJECT)-win

$(NIXBUILD)/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(NIXBUILD)/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

$(NIXBUILD)/libs/%: libs/linux/%
	mkdir -p $(dir $@)
	cp $< $@

$(NIXBUILD)/assets/%: src/assets/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-linux-$(VERSION).zip: \
	$(NIXBUILD)/launch.sh \
	$(NIXBUILD)/$(PROJECT) \
	$(NIXLIBSD) \
	$(subst src/assets/,$(NIXBUILD)/assets/,$(RUNTIME_ASSETS))
	cd build && zip -r $(notdir $@) $(PROJECT)-linux

build: build/$(PROJECT)-osx-$(VERSION).zip

run: build
	$(OSXBUILD)/MacOS/launch.sh
