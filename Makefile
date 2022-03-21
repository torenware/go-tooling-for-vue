export JSASSETS JSMODULE JSDEPS CSSBUNDLE

JSASSETS = dist/assets
JSMODULE = $(wildcard ${JSASSETS}/index.*.js)
JSDEPS = $(wildcard ${JSASSETS}/vendor.*.js)
CSSBUNDLE = $(wildcard ${JSASSETS}/index.*.css)

jsdist: cmd/web/dist/assets cmd/web/main.go
	@echo ${JSMODULE}
	@echo ${JSDEPS}
	@go run ./cmd/web

