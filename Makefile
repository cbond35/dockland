BIN=dockland

all: build

clean:
	go clean

build: clean
	go build -o $(BIN)

