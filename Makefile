SRC_DIR := $(shell ls -d */|grep -vE 'vendor|script|tmp')

all: build

deps:
	# install deps
	@hash dep > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/golang/dep/cmd/dep; \
	fi

	@hash rice > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/GeertJohan/go.rice/rice; \
	fi

	@dep ensure -v

fmt:
	# gofmt code
	gofmt -s -l -w $(SRC_DIR) *.go

rice:
	@rice embed-go

build: rice
	@go build

install: rice
	@go install

test:
	go test -coverprofile .cover.out -v ./...
	# cover
	go tool cover -func=.cover.out
	go tool cover -html=.cover.out -o .cover.html

clean:
	@-rm -f generatemodel
	@rice clean

.PHONY: all deps fmt build test rice install clean