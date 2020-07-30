BIN=dockland
BUILD=go build
CLEAN=go clean

export GO111MODULE=on

all: build

clean:
	$(CLEAN)

build: clean
	$(BUILD) -o $(BIN)

