.PHONY: build install clean run

build:
	go build -o geno-cli ./cmd/geno-cli

install: build
	mkdir -p ~/.geno/bin
	ln -sfn $(PWD)/geno-cli ~/.geno/bin/geno-cli

clean:
	rm -f geno-cli

run: build
	./geno-cli
