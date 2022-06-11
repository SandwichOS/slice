all: build build/slice.slicepkg

build:
	mkdir build

build/slice.slicepkg:
	go build -o build/package/usr/bin/slice
	cp metadata.json build/package
	./build/package/usr/bin/slice create build/package build/slice.slicepkg
clean:
	rm -r build