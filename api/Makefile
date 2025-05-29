all: bin

.PHONY: swag
swag:
	swag fmt
	swag init -g pkg/api/doc.go

.PHONY: jsclient
jsclient:
	if ! [ -d ../ui/src/gen ]; then mkdir ../ui/src/gen; fi;
	rm -rf ../ui/src/gen/apiclient
	docker run --rm -v $(shell pwd):/local -v $(shell pwd)/../ui:/ui --user $(shell id -u):$(shell id -g) \
		openapitools/openapi-generator-cli \
		generate \
		-i /local/docs/swagger.json \
		-g javascript \
		-o /ui/src/gen/apiclient \
		--additional-properties=useES6=true,usePromises=true

.PHONY: ui
ui:
	if ! [ -d pkg/embeded/ui ]; then mkdir pkg/embeded/ui; fi;
	cd ../ui; npm run build
	rm -rfv pkg/embeded/ui/dist
	cp -r ../ui/dist pkg/embeded/ui

.PHONY: bindir
bindir:
	if ! [ -d bin ]; then mkdir bin; fi;

.PHONY: bindir bin
bin: swag jsclient ui
	go build -o bin/go-vue	
