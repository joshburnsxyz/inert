GOX := $(shell which go)
BIN := inert
SRC := cmd/inert

inert:
	mkdir -p dist
	$(GOX) mod tidy
	$(GOX) build \
		-x \
		-v \
		-o ./dist/$(BIN) \
		./$(SRC)

clean:
	@rm -rf dist
	@rm -f *.pem

ssl_keys:
	mkdir -p dist
	$(GOX) mod tidy
	$(GOX) run ./cmd/generate_cert --host 0.0.0.0
