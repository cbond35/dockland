BIN=dockland
BUILD=go build
CLEAN=go clean
TEST=go test -cover

export GO111MODULE=on

all: test build

test:
	$(TEST) ./... -v

clean:
	$(CLEAN)

build: clean
	$(BUILD) -o $(BIN)
