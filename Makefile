default: build

.PHONY: clean
clean:
	rm -rf headers
	rm -rf linux

.PHONY: test
test: headers
	CGO_ENABLED=1 CGO_CFLAGS="-I$$(pwd)/headers/include" go test ./tracer ./printer ./filter

linux:
	git clone --depth 1 https://github.com/torvalds/linux.git ./linux

headers: linux
	cd linux && make headers_install ARCH=x86_64 INSTALL_HDR_PATH=../headers
	rm -rf linux

.PHONY: build
build: headers
	CGO_ENABLED=1 CGO_CFLAGS="-I$$(pwd)/headers/include" go build --ldflags '-linkmode external -extldflags "-static"'

.PHONY: install
install: headers
	CGO_ENABLED=1 CGO_CFLAGS="-I$$(pwd)/headers/include"  go install --ldflags '-linkmode external -extldflags "-static"'

.PHONY: demo
demo: build
	./grace -- cat /dev/null
