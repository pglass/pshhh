#ifndef LEX_H
#define LEX_H

#include <stdlib.h>
#include <string.h>
#include "../buf.h"
#include "../list.h"
#include "textutil.h"
#include "token.h"

typedef struct {
    Buf* buf;
    int i;
    List* tokens;
} Lexer;

Lexer* lexer_init(Buf* buf);
void lexer_free(Lexer* lexer);
int lexer_next(Lexer* lexer);
int lexer_peek(Lexer* lexer);
void lexer_lex(Lexer* lexer);
void lexer_read_whitespace(Lexer* lexer);
void lexer_read_alpha_numeric(Lexer* lexer);
void lexer_read_punctuation(Lexer* lexer);
void lexer_yield_token(Lexer* lexer, Buf* buf, TokenType type);
void lexer_yield_token_from_str(Lexer* lexer, char* s, TokenType type);

#endif
