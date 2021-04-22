.PHONY: fmt
fmt:
	gofmt -w src/*.go

.PHONY: unit-test
unit-test: fmt
	go test -v tests/

.PHONY: build
build:
	cd src && go build -o ../stegasaurus

.PHONY: clean
clean:
	rm -f stegasaurus
