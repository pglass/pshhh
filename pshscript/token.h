#ifndef TOKEN_H
#define TOKEN_H

#include <stdlib.h>
#include <stdio.h>
#include "tokentype.h"

typedef union {
    int i;
    float f;
    char* str;
} TokenData;

typedef enum {
    UNION_INT,
    UNION_FLOAT,
    UNION_STRING,
} TokenUnionType;

typedef struct {
    TokenType type;
    TokenUnionType uniontype;
    TokenData data;
} Token;

/* No copy is made of str */
Token* token_init_str(char* str, TokenType type);
void token_print(Token* token);
void token_free(Token* token);

#endif
