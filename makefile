
CC=clang
CFLAGS=-Wall --std=c99 -g

all: pshhh

buf.o: buf.h buf.c
	$(CC) $(CFLAGS) -c $^

pshhh: pshhh.c buf.o
	$(CC) $(CFLAGS) $^ -o $@

run: pshhh
	./pshhh

clean:
	rm -f *.o *.gch pshhh
