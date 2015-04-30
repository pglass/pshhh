#include "lex.h"

#define LEX_TMP_BUF_INITIAL_CAPACITY 40
#define LEX_TMP_STRING_BUF_CAPACITY 128

Lexer* lexer_init(Buf* buf) {
    Lexer* lexer = (Lexer*) malloc(sizeof(Lexer));
    lexer->buf = buf;
    lexer->i = 0;
    lexer->tokens = list_init();
    return lexer;
}

void lexer_free(Lexer* lexer) {
    free(lexer);
}

int lexer_next(Lexer* lexer) {
    int c = lexer_peek(lexer);
    if (c > 0) lexer->i++;
    return c;
}

int lexer_peek(Lexer* lexer) {
    if (lexer->i >= lexer->buf->len) {
        return -1;
    }
    return lexer->buf->buf[lexer->i];
}

void lexer_yield_token(Lexer* lexer, Buf* buf, TokenType type) {
    char* s = (char*) malloc(buf->len);
    strncpy(s, buf->buf, buf->len);
    Token* token = token_init_str(s, type);
    list_append(lexer->tokens, (void*) token);
}

void lexer_yield_token_from_str(Lexer* lexer, char* s, TokenType type) {
    int n = strlen(s);
    char* new = (char*) malloc(n);
    strncpy(new, s, n);
    Token* token = token_init_str(s, type);
    list_append(lexer->tokens, (void*) token);
}

/*
 * Input:   lexer->buf
 * Output:  lexer->tokens
 */
void lexer_lex(Lexer* lexer) {
    int c;
    while ((c = lexer_peek(lexer)) >= 0) {
        if (is_whitespace(c)) {
            lexer_read_whitespace(lexer);
        } else if (is_digit(c) || is_alpha(c) || c == '_') {
            lexer_read_alpha_numeric(lexer);
        } else if (c == '\'') {
            lexer_read_single_quoted_string(lexer);
        } else if (c == '"') {
            lexer_read_double_quoted_string(lexer);
        } else if (is_punctuation(c)) {
            lexer_read_punctuation(lexer);
        } else {
            printf("failed to handle char '%c'\n", (char) c);
            lexer_next(lexer);
        }
    }
}

void lexer_read_whitespace(Lexer* lexer) {
    Buf* buf = buf_init(LEX_TMP_BUF_INITIAL_CAPACITY);
    int c;
    while ((c = lexer_peek(lexer)) >= 0 && is_whitespace(c)) {
        // yield one token for each newline.
        // yield one token for each contiguous section of whitespace.
        if (c == '\n') {
            lexer_yield_token(lexer, buf, WHITESPACE);
            buf_clear(buf);
            buf_append_char(buf, c);
            lexer_yield_token(lexer, buf, NEWLINE);
            buf_clear(buf);
        } else {
            buf_append_char(buf, c);
        }
        lexer_next(lexer);
    }
    if (buf->len > 0) {
        lexer_yield_token(lexer, buf, WHITESPACE);
    }
    buf_free(buf);
}

void handle_reserved_word(Token* token) {
    assert(token != NULL);
         if (!strcmp(token->data.str, "if")) token->type = IF;
    else if (!strcmp(token->data.str, "then")) token->type = THEN;
    else if (!strcmp(token->data.str, "else")) token->type = ELSE;
    else if (!strcmp(token->data.str, "elif")) token->type = ELIF;
    else if (!strcmp(token->data.str, "fi")) token->type = FI;
    else if (!strcmp(token->data.str, "do")) token->type = DO;
    else if (!strcmp(token->data.str, "done")) token->type = DONE;
    else if (!strcmp(token->data.str, "case")) token->type = CASE;
    else if (!strcmp(token->data.str, "esac")) token->type = ESAC;
    else if (!strcmp(token->data.str, "while")) token->type = WHILE;
    else if (!strcmp(token->data.str, "until")) token->type = UNTIL;
    else if (!strcmp(token->data.str, "for")) token->type = FOR;
    else if (!strcmp(token->data.str, "in")) token->type = IN;
    else if (!strcmp(token->data.str, "function")) token->type = FUNCTION;
    else if (!strcmp(token->data.str, "select")) token->type = SELECT;
}

void lexer_read_alpha_numeric(Lexer* lexer) {
    Buf* buf = buf_init(LEX_TMP_BUF_INITIAL_CAPACITY);
    int c;
    while ((c = lexer_peek(lexer)) >= 0
            && (is_alpha(c) || is_digit(c) || c == '_')) {
        buf_append_char(buf, c);
        lexer_next(lexer);
    }
    /* A WORD that does not start with a digit is a NAME */
    if (buf->len > 0 && !is_digit(buf->buf[0])) {
        lexer_yield_token(lexer, buf, NAME);
        handle_reserved_word((Token*) lexer->tokens->tail->prev->item);
    } else if (buf->len > 0) {
        lexer_yield_token(lexer, buf, WORD);
    }
    buf_free(buf);
}

void lexer_read_punctuation(Lexer* lexer) {
    int c = lexer_next(lexer);
    int cc = lexer_peek(lexer);
    char s[3] = "\0\0";
    s[0] = c;
    if (c >= 0 && cc >= 0) s[1] = cc;

    // important to check greedily here. longest match first.
    // e.g. parse one "<<" token instead of two "<" tokens.

    if (!strcmp(s, "<<")) {
        lexer_next(lexer);
        if (lexer_peek(lexer) == '-') {
            lexer_yield_token_from_str(lexer, "<<-", DLESSDASH);
            lexer_next(lexer);
        } else {
            lexer_yield_token_from_str(lexer, "<<", DLESS);
        }
    } else if (!strcmp(s, "&&")) {
        lexer_yield_token_from_str(lexer, "&&", AND_IF); lexer_next(lexer);
    } else if (!strcmp(s, "||")) {
        lexer_yield_token_from_str(lexer, "||", OR_IF); lexer_next(lexer);
    } else if (!strcmp(s, ";;")) {
        lexer_yield_token_from_str(lexer, ";;", DSEMI); lexer_next(lexer);
    } else if (!strcmp(s, ">>"))  {
        lexer_yield_token_from_str(lexer, ">>", DGREAT); lexer_next(lexer);
    } else if (!strcmp(s, "<&")) {
        lexer_yield_token_from_str(lexer, "<&", LESSAND); lexer_next(lexer);
    } else if (!strcmp(s, ">&")) {
        lexer_yield_token_from_str(lexer, ">&", GREATAND); lexer_next(lexer);
    } else if (!strcmp(s, "<>")) {
        lexer_yield_token_from_str(lexer, "<>", LESSGREAT); lexer_next(lexer);
    } else if (!strcmp(s, ">|")) {
        lexer_yield_token_from_str(lexer, ">|", CLOBBER); lexer_next(lexer);
    } else if (!strcmp(s, "[[")) {
        lexer_yield_token_from_str(lexer, "[[", DLBRACKET); lexer_next(lexer);
    } else if (!strcmp(s, "]]")) {
        lexer_yield_token_from_str(lexer, "]]", DRBRACKET); lexer_next(lexer);
    } else switch (c) {
        case '|': lexer_yield_token_from_str(lexer, "|", PIPE); break;
        case '&': lexer_yield_token_from_str(lexer, "&", AMPERSAND); break;
        case ';': lexer_yield_token_from_str(lexer, ";", SEMI); break;
        case '<': lexer_yield_token_from_str(lexer, "<", LESS); break;
        case '>': lexer_yield_token_from_str(lexer, ">", GREATER); break;
        case '(': lexer_yield_token_from_str(lexer, "(", LPAREN); break;
        case ')': lexer_yield_token_from_str(lexer, ")", RPAREN); break;
        case '$': lexer_yield_token_from_str(lexer, "$", DOLLAR); break;
        case '`': lexer_yield_token_from_str(lexer, "`", BACKTICK); break;
        case '\\': lexer_yield_token_from_str(lexer, "\\", BACKSLASH); break;
//        case '\'': lexer_yield_token_from_str(lexer, "\'", QUOTE); break;
//        case '"': lexer_yield_token_from_str(lexer, "\"", DQUOTE); break;
        case '+': lexer_yield_token_from_str(lexer, "+", PLUS); break;
        case '-': lexer_yield_token_from_str(lexer, "-", DASH); break;
        case '*': lexer_yield_token_from_str(lexer, "*", ASTERISK); break;
        case '/': lexer_yield_token_from_str(lexer, "/", SLASH); break;
        case '?': lexer_yield_token_from_str(lexer, "?", QUESTION); break;
        case '[': lexer_yield_token_from_str(lexer, "[", LBRACKET); break;
        case ']': lexer_yield_token_from_str(lexer, "]", RBRACKET); break;
        case '#': lexer_yield_token_from_str(lexer, "#", HASH); break;
        case '~': lexer_yield_token_from_str(lexer, "~", TILDE); break;
        case '=': lexer_yield_token_from_str(lexer, "=", EQUALS); break;
        case '%': lexer_yield_token_from_str(lexer, "%", PERCENT); break;
        case '{': lexer_yield_token_from_str(lexer, "{", LBRACE); break;
        case '}': lexer_yield_token_from_str(lexer, "}", RBRACE); break;
        case '!': lexer_yield_token_from_str(lexer, "!", BANG); break;
        default:
            printf("failed to handle character '%c'\n", c);
    }
}

/* We have options for tokenizing double-quoted strings. If our input is,
 *      "var=$(VAR) aaa "
 * then we could produce a single token,
 *      ("var is $(VAR)" STRING)
 * But this forces us to retokenize later. It would be easier to tokenize a bit
 * now, in a way that gives enough information to do expansions while still
 * being able to reconstruct the original string.
 *      ('"' DQUOTE)
 *      ("var is " STR_SEGMENT)
 *      ("$" DOLLAR)
 *      ("(" LPAREN)
 *      ("VAR" NAME)
 *      (")" RPAREN)
 *      (" aaa " STR_SEGMENT)
 *      ('"' DQUOTE)
 * We have to be careful not to lose any portion of the original string:
 *    - Single-quoted strings will be parsed as a single STRING token, with no
 *    - The two delimiting double-quote chars produce a DQUOTE token.
 *    - Expansions in double-quoted strings will be tokenized fully.
 *    - Each segment of contiguous text not part of an expansion in a
 *      double-quoted string will yield a single STR_SEGMENT token.
 */
void lexer_read_single_quoted_string(Lexer* lexer) {
    Buf* buf = buf_init(LEX_TMP_STRING_BUF_CAPACITY);
    // don't put the single quotes in the buf.
    int c = lexer_next(lexer);
    assert(c == '\'');
    while ((c = lexer_peek(lexer)) >= 0 && c != '\'') {
        buf_append_char(buf, lexer_next(lexer));
    }
    if (c != '\'') {
        // TODO: if unclosed, tell the repl to get multline input and retry
        printf("ERROR: unclosed string '%s", buf->buf);
    } else {
        lexer_yield_token(lexer, buf, SINGLE_QUOTED_STRING);
        lexer_next(lexer);
    }
    buf_free(buf);
}

void lexer_read_double_quoted_string(Lexer* lexer) {
    Buf* buf = buf_init(LEX_TMP_STRING_BUF_CAPACITY);
    int c = lexer_next(lexer);
    assert(c == '"');
    lexer_yield_token_from_str(lexer, "\"", DQUOTE);
    while ((c = lexer_peek(lexer)) >= 0 && c != '"') {
        // convert backslash escapes
        if (c == '\\') {
            lexer_next(lexer);
            c = lexer_peek(lexer);
            if (is_char(c, "$`\"\\")) {
                buf_append_char(buf, c);
                lexer_next(lexer);
            } else if (c == '\n') {
                // escaped newlines are ignored
            } else {
                buf_append_char(buf, '\\');
            }
        } else if (c == '$') {
            // terminate and yield the current string segment token.
            // yield tokens for an expansion. resume reading a new
            // string segment
            lexer_yield_token(lexer, buf, STR_SEGMENT);
            buf_clear(buf);
            lexer_read_dollar_expansion(lexer);
            buf_clear(buf);
        } else {
            buf_append_char(buf, lexer_next(lexer));
        }
    }
    if (c != '"') {
        // TODO: have the parser recognize the unclosed string an dproduce a
        // better error message
        printf("ERROR: Unclosed string\n");
    } else {
        assert(c == '"');
        if (buf->len > 0) {
            lexer_yield_token(lexer, buf, STR_SEGMENT);
        }
        lexer_yield_token_from_str(lexer, "\"", DQUOTE);
        lexer_next(lexer);
    }
}


void lexer_read_dollar_expansion(Lexer* lexer) {
    int c = lexer_next(lexer);
    assert(c == '$');
    lexer_yield_token_from_str(lexer, "$", DOLLAR);

    c = lexer_peek(lexer);
    if (c == '(') {
        lexer_read_expansion_parens(lexer);
    } else if (c == '{') {
        lexer_read_expansion_braces(lexer);
    } else if (is_alpha(c) || c == '_') {
        lexer_read_alpha_numeric(lexer);
    } else if (c == 0) {
        printf("ERROR: unexpected end of input\n");
    } else {
        printf("ERROR: invalid expansion sequence '$%c'\n", c);
    }
}

// Read the '()' whenever we an expansion like '$()'
// The expansions are allowed to be recursive
void lexer_read_expansion_parens(Lexer* lexer) {
    int c = lexer_next(lexer);
    assert(c == '(');
    lexer_yield_token_from_str(lexer, "(", LPAREN);
    int depth = 1;
    while (depth > 0 && (c = lexer_peek(lexer)) >= 0) {
        if (c == '$') {
            lexer_read_dollar_expansion(lexer);
        } else if (c == ')') {
            lexer_yield_token_from_str(lexer, ")", RPAREN);
            lexer_next(lexer);
            depth--;
        } else if (c == '(') {
            lexer_yield_token_from_str(lexer, "(", LPAREN);
            lexer_next(lexer);
            depth++;
        } else if (is_digit(c) || is_alpha(c) || c == '_') {
            lexer_read_alpha_numeric(lexer);
        } else {
            printf("ERROR: unhandled char in parens expansion: '%c'\n", c);
        }
    }
}


void lexer_read_expansion_braces(Lexer* lexer) {
    Buf* buf = buf_init(LEX_TMP_BUF_INITIAL_CAPACITY);
    int c = lexer_next(lexer);
    assert(c == '{');
    lexer_yield_token_from_str(lexer, "{", LBRACE);
    int found_split = 0;

    while ((c = lexer_peek(lexer)) >= 0 && c != '}') {
        // handle an escaped \{
        if (c == '\\') {
            lexer_next(lexer);
            if (lexer_peek(lexer) == '{') {
                buf_append_char(buf, '{');
                lexer_next(lexer);
            } else {
                buf_append_char(buf, '\\');
            }
        // handle special operators like the ":-" in ${param:-foo}
        // these include + - = ? :+ :- := :?
        } else if (!found_split && is_char(c, ":+-=?")) {
            found_split = 1;
            if (buf->len > 0) {
                lexer_yield_token(lexer, buf, STR_SEGMENT);
                buf_clear(buf);
            } else {
                printf("ERROR: no parameter found in parameter expansion\n");
            }
            if (c == ':') {
                lexer_yield_token_from_str(lexer, ":", COLON);
                lexer_next(lexer);
            }
            c = lexer_peek(lexer);
            switch (c) {
                case '+': lexer_yield_token_from_str(lexer, "+", PLUS);
                          lexer_next(lexer);
                          break;
                case '-': lexer_yield_token_from_str(lexer, "-", DASH);
                          lexer_next(lexer);
                          break;
                case '=': lexer_yield_token_from_str(lexer, "=", EQUALS);
                          lexer_next(lexer);
                          break;
                case '?': lexer_yield_token_from_str(lexer, "?", QUESTION);
                          lexer_next(lexer);
                          break;
            }
        // handle other special operators
        //   string length ${#foo}
        //   remove suffixes ${param%word} ${param%%word}
        //   remove prefixes ${param#word} ${param##word}
        } else if (!found_split && is_char(c, "%#")) {
            found_split = 1;
            if (buf->len > 0) {
                lexer_yield_token(lexer, buf, STR_SEGMENT);
                // buf is cleared later
            }
            if (c == '#') {
                lexer_yield_token_from_str(lexer, "#", HASH);
                lexer_next(lexer);
                // a double '##' is only allowed with a non-empty parameter
                if (buf->len > 0 && lexer_peek(lexer) == '#') {
                    lexer_yield_token_from_str(lexer, "#", HASH);
                    lexer_next(lexer);
                }
            } else if (c == '%') {
                lexer_yield_token_from_str(lexer, "%", PERCENT);
                lexer_next(lexer);
                if (lexer_peek(lexer) == '%') {
                    lexer_yield_token_from_str(lexer, "%", PERCENT);
                    lexer_next(lexer);
                }
            }
            buf_clear(buf);
        } else {
            buf_append_char(buf, c);
            lexer_next(lexer);
        }
    }
    if (c == '}') {
        lexer_yield_token(lexer, buf, STR_SEGMENT);
        lexer_yield_token_from_str(lexer, "}", RBRACE);
        lexer_next(lexer);
    } else {
        printf("ERROR: Unclosed parameter expansion ${\n");
    }
    buf_free(buf);
}
