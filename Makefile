export JSASSETS

JSASSETS = dist/assets

cmd/web/dist/assets: $(shell find src -type f -name '*.vue' -o -name '*.js' -o -name '*.ts')
	@yarn build

go: cmd/web/dist/assets cmd/web/main.go
	@echo ${JSASSETS}
	@go run ./cmd/web

rebuild:  go
