.PHONY: test

main: main.go lex/*.go ast/*.go exe/*.go
	go build $<

test:
	go test ./test/...

clean:
	rm -f main
