SHELL = /bin/bash

GO := go
DOCKER := docker
TAILWIND_VERSION := "v4.1.7"
TEMPL := templ
MINIFY := minify

gofiles != find src/webserver/ -name '*.go' -type f
templgo != find src/webserver/templ_components/ -name '*.go' -type f
buildtemplgo := $(templgo:src/webserver/templ_components/%.go=src/webserver/components/%.go)
sql != find src/sql/* -type f 
buildsql := $(sql:src/sql/%=build/sql/%)
fonts != find src/fonts/*.ttf -type f
buildfonts := $(fonts:src/fonts/%.ttf=build/fonts/%.ttf)
js != find src/js/*.js -type f
templfiles != find src/webserver/templ_components/ -name '*.templ' -type f
buildtempl := $(templfiles:src/webserver/templ_components/%.templ=src/webserver/components/%_templ.go)
icons != find src/icons/* -type f
buildicons := $(icons:src/icons/%=build/icons/%)

.PHONY: clean build run templ tailwindcss configure complete docker

build: build/css/stylesheet.css build/cmd/server $(buildsql) $(buildfonts) build/js/out.js $(buildicons) build/tls/cert.pem build/tls/key.pem 

# Cleans downloaded and generated files. Changing permissions in tools in case tools is also used for GOPATH
clean:
	rm -rf build/
	chmod -R 777 tools/ || true
	rm -rf tools/ 
	find src/webserver/components/ -name '*_templ.go' -type f -exec rm '{}' \;

run:
	$(MAKE) build
	cd build/cmd && ./server

templ:
	make $(buildtempl)
	make $(buildtemplgo)

tailwindcss:
	$(MAKE) build/css/stylesheet.css

# Generates self-signed certificates for build if none exist already. This presumes Go is installed in the location suggested by the official Go docs.
# If certificates need to be renewed for testing, simply delete the files.
# In production, use proper certificates signed by an authority rather than the ones generated here.
build/tls/cert.pem build/tls/key.pem: | build/cmd/server
	mkdir -p $$(dirname $@)	
	cd build/tls/ && $(GO) run $(shell dirname $(shell type -P go))/../src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=127.0.0.1

# Generates Tailwindcss if Templ HTML has been updated, if the the input.css file has been updated, or if the javascript files have been updated.
# Any of these three files could result in new TailwindCSS utility classes being added.
# Removes old stylesheet.css in the case that the prerequistes were updated but no new utility classes were called resulting in the stylesheet.css being skipped by the
# TailwindCSS cli tool
build/css/stylesheet.css : $(buildtempl) src/css/input.css build/js/out.js | tools/tailwindcss 
	rm -f build/css/stylesheet.css
	tools/tailwindcss -m  -i src/css/input.css -o build/css/stylesheet.css 

$(buildtempl):src/webserver/components/%_templ.go:src/webserver/templ_components/%.templ
	mkdir -p $$(dirname $@)
	cd src/webserver/templ_components && $(TEMPL) generate -f $*.templ
	mv src/webserver/templ_components/$*_templ.go $$(dirname $@)/

$(buildtemplgo):src/webserver/components/%.go:src/webserver/templ_components/%.go
	mkdir -p $$(dirname $@)
	cp src/webserver/templ_components/$*.go src/webserver/components/$*.go

# Generates built executable of webserver
build/cmd/server : $(gofiles) $(buildtempl) $(buildtemplgo)
	$(GO) build -C src/webserver/main -o ../../../build/cmd/server && chmod +x build/cmd/server

# Copy any updated sql files
$(buildsql):build/sql/% :src/sql/%
	mkdir -p $$(dirname $@)
	cp $< $@

# Copy font files to build directory
$(buildfonts):build/fonts/%.ttf:src/fonts/%.ttf
	mkdir -p $$(dirname $@)
	cp $< $@

# Minify Javascript scripts to build directory
build/js/out.js: $(js)
	mkdir -p build/js
	$(MINIFY) -b -r -o ./build/js/out.js ./src/js/

# Copy icons to build directory
$(buildicons):build/icons/%:src/icons/%
	mkdir -p $$(dirname $@)
	cp $< $@

configure:
#	Install project dependencies
	cd src/webserver/ && $(GO) mod download
#	Install minify binary to GOPATH
	$(GO) install github.com/tdewolff/minify/v2/cmd/minify@v2.23.5
# 	Install Templ
	$(GO) install github.com/a-h/templ/cmd/templ@v0.3.857
# Install TailwindCSS-CLI standalone executable
tools/tailwindcss:
	mkdir -p tools
	wget -O tools/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-linux-x64
	chmod +x tools/tailwindcss

# Handles dependency configuration and building
complete:
	make configure
	make build

docker:
	sudo $(DOCKER) build -t com.ethanrandolph.webserver .
