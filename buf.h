#ifndef BUF_H
#define BUF_H

#include <stdlib.h>
#include <assert.h>
#include <string.h>

typedef struct {
    char* buf;
    size_t capacity;
    size_t len;
} Buf;

Buf* buf_init(size_t initial_capacity);
void buf_free(Buf* buf);
void buf_append(Buf* buf, char* data, size_t len);
void buf_append_char(Buf* buf, char c);
void buf_strip_last_char(Buf* buf);
char buf_get_last_char(Buf* buf);
void buf_clear(Buf* buf);
char* buf_to_str(Buf* buf);

#endif
