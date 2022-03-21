export JSASSETS JSMODULE JSDEPS CSSBUNDLE

JSASSETS = dist/assets
JSMODULE = $(wildcard ${JSASSETS}/index.*.js)
JSDEPS = $(wildcard ${JSASSETS}/vendor.*.js)
CSSBUNDLE = $(wildcard ${JSASSETS}/index.*.css)

cmd/web/dist/assets: $(shell find src -type f -name '*.vue' -o -name '*.js' -o -name '*.ts')
	@yarn build

go: cmd/web/dist/assets cmd/web/main.go
	@echo ${JSMODULE}
	@echo ${JSDEPS}
	@go run ./cmd/web

rebuild:  go
