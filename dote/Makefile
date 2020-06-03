

.PHONY: build testex

main-file:
	toml2cli --in-file=clifile.toml --out-file=main.go

build: main-file
	go build -o dote .

testex: build
	./dote get --source github.com/ttacon/dotefiles --profile=dev

testdiag: build
	./dote diag
