.PHONY: build install clean run

build:
	go build -o geno ./cmd/geno

install: build
	cp geno /usr/local/bin/geno

clean:
	rm -f geno

run: build
	./geno
