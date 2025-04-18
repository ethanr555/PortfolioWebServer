go != find src/webserver/*.go
sql != find src/sql/* -type f 
buildsql := $(sql:src/sql/%=build/sql/%)
fonts != find src/fonts/*.ttf -type f
buildfonts := $(fonts:src/fonts/%.ttf=build/fonts/%.ttf)
js != find src/js/*.js -type f
buildjs := $(js:src/js/%.js=build/js/%.js)
templ != find src/webserver/*.templ
buildtempl := $(templ:%.templ=%_templ.go)

.PHONY: clean build run templ tailwindcss

build: build/css/stylesheet.css build/bin/server $(buildsql) $(buildfonts) $(buildjs) build/tls/cert.pem build/tls/key.pem

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

build/tls/cert.pem build/tls/key.pem: | build/bin/server
	mkdir -p $$(dirname $@)	
	cd build/tls/ && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

# Generates Tailwindcss if html files have been updated or the input.css file was updated.
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

$(buildfonts):build/fonts/%.ttf:src/fonts/%.ttf
	mkdir -p $$(dirname $@)
	cp $< $@

$(buildjs):build/js/%.js:src/js/%.js
	mkdir -p $$(dirname $@)
	cp $< $@