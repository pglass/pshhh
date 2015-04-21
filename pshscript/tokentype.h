#ifndef TOKENTYPE_H
#define TOKENTYPE_H
#include <assert.h>

/* DON'T ADD THINGS HERE WITHOUT CAREFULLY UPDATING token_type_to_str */
typedef enum {
    WORD,
    ASSIGNMENT_WORD,
    NAME,
    NEWLINE,
    WHITESPACE, // excluding '\n'
    IO_NUMBER,

    MIN_PUNCTUATION,
    PIPE,       // |
    AMPERSAND,  // &
    SEMI,       // ;
    LESS,       // <
    GREATER,    // >
    LPAREN,     // (
    RPAREN,     // )
    DOLLAR,     // $
    BACKTICK,   // `
    BACKSLASH,  /* \ */
    QUOTE,      // '
    DQUOTE,     // "
    PLUS,       // +
    DASH,       // -
    ASTERISK,   // *
    SLASH,      // /
    QUESTION,   // ?
    LBRACKET,   // [
    RBRACKET,   // ]
    HASH,       // #
    TILDE,      // ~
    EQUALS,     // =
    PERCENT,    // %
    LBRACE,     // {
    RBRACE,     // }
    BANG,       // !

    AND_IF,     // &&
    OR_IF,      // ||
    DSEMI,      // ;;
    DLESS,      // <<
    DGREAT,     // >>
    LESSAND,    // <&
    GREATAND,   // >&
    LESSGREAT,  // <>
    DLESSDASH,  // <<-
    CLOBBER,    // >|
    DLBRACKET,  // [[
    DRBRACKET,  // ]]
    MAX_PUNCTUATION,

    MIN_RESERVED_WORD,
    IF,
    THEN,
    ELSE,
    ELIF,
    FI,
    DO,
    DONE,
    CASE,
    ESAC,
    WHILE,
    UNTIL,
    FOR,
    IN,
    FUNCTION,
    SELECT,
    MAX_RESERVED_WORD,
    N_TOKEN_TYPES,
} TokenType;

char* tokentype_to_str(TokenType type);

#endif
