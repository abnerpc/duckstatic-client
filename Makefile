
APP_NAME=ds-client
DEBUG_APP_NAME=ds-client-debug

.PHONY: clean
clean:
	@rm -f $(APP_NAME)

.PHONY: build
build: clean
	@go build -o $(APP_NAME)

build-debug: clean
	@godebug build -o $(DEBUG_APP_NAME)
