CC=clang
CFLAGS=-Wall --std=c99 -g

PSHSCRIPT_DIR=./pshscript

all: pshhh

buf.o: buf.h buf.c
	$(CC) $(CFLAGS) -c $^

list.o: list.h list.c
	$(CC) $(CFLAGS) -c $^

lex.o: $(PSHSCRIPT_DIR)/lex.h $(PSHSCRIPT_DIR)/lex.c
	$(CC) $(CFLAGS) -c $^

token.o: $(PSHSCRIPT_DIR)/token.h $(PSHSCRIPT_DIR)/token.c
	$(CC) $(CFLAGS) -c $^

tokentype.o: $(PSHSCRIPT_DIR)/tokentype.h $(PSHSCRIPT_DIR)/tokentype.c
	$(CC) $(CFLAGS) -c $^

textutil.o: $(PSHSCRIPT_DIR)/textutil.h $(PSHSCRIPT_DIR)/textutil.c
	$(CC) $(CFLAGS) -c $^

pshhh: pshhh.c buf.o list.o lex.o textutil.o token.o tokentype.o
	$(CC) $(CFLAGS) $^ -o $@

run: pshhh
	./pshhh

clean:
	find . -name "*.o" -or -name "*.gch" -or -name "*.dSYM" -delete
	rm -f pshhh lex_test_shim

#=============================
# TEST RULES
#=============================
TEST_DIR=./test

lex_test_shim: $(TEST_DIR)/lex_test_shim.c lex.o token.o tokentype.o textutil.o list.o buf.o
	$(CC) $(CFLAGS) $^ -o $@

run_tests: lex_test_shim
	python $(TEST_DIR)/test_lex.py "./lex_test_shim"
