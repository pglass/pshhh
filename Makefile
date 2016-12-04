.PHONY: test

main: main.go lex/*.go ast/*.go exe/*.go
	go build $<

test:
	# go test ./test/...
	go test ./test/lex_test.go

clean:
	rm -f main
