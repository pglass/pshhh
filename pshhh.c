#include <stdio.h>
#include "buf.h"

#define BUF_SIZE 1000

void prompt(char* s) {
    printf("%s ", s);
}

void read_line(Buf* buf) {
    buf_clear(buf);
    char data[BUF_SIZE];
    while (1) {
        memset(data, 0, BUF_SIZE);
        if (fgets(data, BUF_SIZE, stdin) != NULL) {
            buf_append(buf, data, strlen(data));
            if (buf_get_last_char(buf) == '\n') {
                buf_strip_last_char(buf);
                return;
            }
        }
    }
}

int main(int argc, char** argv) {
    printf("Starting pshhh\n");
    Buf* buf = buf_init(BUF_SIZE);
    prompt("$");
    read_line(buf);
    printf("Got input: %s\n", buf->buf);
}
