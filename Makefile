go != find src/webserver/*.go
sql != find src/sql/* -type f 
buildsql := $(sql:src/sql/%=build/sql/%)
fonts != find src/fonts/*.ttf -type f
buildfonts := $(fonts:src/fonts/%.ttf=build/fonts/%.ttf)
js != find src/js/*.js -type f
buildjs := $(js:src/js/%.js=build/js/%.js)
templ != find src/webserver/*.templ
buildtempl := $(templ:%.templ=%_templ.go)
icons != find src/icons/* -type f
buildicons := $(icons:src/icons/%=build/icons/%)

.PHONY: clean build run templ tailwindcss

build: build/css/stylesheet.css build/bin/server $(buildsql) $(buildfonts) $(buildjs) $(buildicons) build/tls/cert.pem build/tls/key.pem

clean:
	rm -rf build/
	rm src/webserver/*_templ.go

run:
	$(MAKE) build
	cd build/bin/ && ./server

templ:
	templ generate -path src/webserver/

tailwindcss:
	$(MAKE) build/css/stylesheet.css

# Generates self-signed certificates for build if none exist already. This presumes Go is installed in the location suggested by the official Go docs.
# If certificates need to be renewed for testing, simply delete the files.
# In production, use proper certificates signed by an authority rather than the ones generated here.
build/tls/cert.pem build/tls/key.pem: | build/bin/server
	mkdir -p $$(dirname $@)	
	cd build/tls/ && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

# Generates Tailwindcss if Templ HTML has been updated, if the the input.css file has been updated, or if the javascript files have been updated.
# Any of these three files could result in new TailwindCSS utility classes being added.
# Removes old stylesheet.css in the case that the prerequistes were updated but no new utility classes were called resulting in the stylesheet.css being skipped by the
# TailwindCSS cli tool
build/css/stylesheet.css : $(buildtempl) src/css/input.css $(buildjs)
	rm -f build/css/stylesheet.css
	npx @tailwindcss/cli -m  -i src/css/input.css -o build/css/stylesheet.css 

$(buildtempl):src/webserver/%_templ.go:src/webserver/%.templ
	cd src/webserver && templ generate -f $*.templ

# Generates built executable of webserver
build/bin/server : $(go) $(buildtempl)
	echo $(buildhtml)
	go build -C src/webserver/ -o ../../build/bin/server

# Copy any updated sql files
$(buildsql):build/sql/% :src/sql/%
	mkdir -p $$(dirname $@)
	cp $< $@

# Copy font files to build directory
$(buildfonts):build/fonts/%.ttf:src/fonts/%.ttf
	mkdir -p $$(dirname $@)
	cp $< $@

# Copy Javascript scripts to build directory
$(buildjs):build/js/%.js:src/js/%.js
	mkdir -p $$(dirname $@)
	cp $< $@

# Copy icons to build directory
$(buildicons):build/icons/%:src/icons/%
	mkdir -p $$(dirname $@)
	cp $< $@
