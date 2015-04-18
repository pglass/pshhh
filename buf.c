#include "buf.h"
#include <stdio.h>

Buf* buf_init(size_t initial_capacity) {
    Buf* buf = (Buf*) malloc(sizeof(Buf));
    buf->buf = (char*) malloc(initial_capacity * sizeof(char));
    buf->capacity = initial_capacity;
    buf->len = 0;
    return buf;
}

void buf_free(Buf* buf) {
    if (buf != NULL) {
        return;
    }
    free(buf->buf);
    free(buf);
}

/* Guarantees buf's capacity will be at least capacity */
void buf_ensure_capacity(Buf* buf, size_t capacity) {
    assert(buf != NULL);
    if (buf->capacity >= capacity) {
        return;
    }
    buf->capacity = capacity * 1.5;
    buf->buf = realloc(buf->buf, buf->capacity * sizeof(char));
    assert(buf->capacity >= buf->len);
}

void buf_append(Buf* buf, char* data, size_t len) {
    assert(buf != NULL);
    assert(buf->buf != NULL);
    size_t new_len = buf->len + len;
    buf_ensure_capacity(buf, new_len);
    memcpy(buf->buf + buf->len, data, len);
    buf->len = new_len;
}

char buf_get_last_char(Buf* buf) {
    assert(buf != NULL);
    return (buf->len == 0) ? 0 : buf->buf[buf->len - 1];
}

void buf_strip_last_char(Buf* buf) {
    assert(buf != NULL);
    if (buf->len > 0) {
        buf->len -= 1;
        buf->buf[buf->len] = 0;
    }
}

void buf_clear(Buf* buf) {
    memset(buf->buf, 0, buf->len);
    buf->len = 0;
}
