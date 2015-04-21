#include "token.h"

/* No copy is made of str */
Token* token_init_str(char* str, TokenType type) {
    Token* token = (Token*) malloc(sizeof(Token));
    token->data.str = str;
    token->uniontype = UNION_STRING;
    token->type = type;
    return token;
}

void token_print(Token* token) {
    printf("Token[data='%s' type=%s]\n",
        token->data.str, tokentype_to_str(token->type));
}

void token_free(Token* token) {
    free(token);
}


