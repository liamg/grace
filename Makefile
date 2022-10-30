default: build

.PHONY: test
test:
	CC=musl-gcc CGO_ENABLED=1 go test ./...

.PHONY: build
build:
	CGO_ENABLED=1 CC=musl-gcc go build --ldflags '-linkmode external -extldflags "-static"'

.PHONY: demo
demo: build
	./grace -- cat /dev/null
