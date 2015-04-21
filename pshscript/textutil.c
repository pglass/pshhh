#include "textutil.h"

/* Check if c matches any of the given chars
 * returns 1 on a match. 0 on failure.
 */
int is_char(int c, char* choices) {
    size_t len = strlen(choices);
    for (int i = 0; i < len; ++i) {
        if (c == choices[i]) {
            return 1;
        }
    }
    return 0;
}

int is_whitespace(int c) {
    return is_char(c, " \t\n\r");
}

int is_alpha(int c) {
    return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z');
}

int is_digit(int c) {
    return '0' <= c && c <= '9';
}

int is_punctuation(int c) {
    // no underscore
    return is_char(c, "`~!@#$%^&*()+-={}[]|\\:;\"'<>,.?/");
}
