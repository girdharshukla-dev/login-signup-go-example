.PHONY: build run clean

build:
	mkdir -p build
	go build -o build/login_service

run: build
	./build/login_service

clean:
	rm -rf build/

