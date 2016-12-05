.PHONY: test

psh: psh.go lex/*.go ast/*.go exe/*.go
	go build $<

test: psh
	go test ./test/...

clean:
	rm -f psh
