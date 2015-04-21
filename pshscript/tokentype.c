#include "tokentype.h"

char* TTYPE_TO_STR[] = {
    "WORD",
    "ASSIGNMENT_WORD",
    "NAME",
    "NEWLINE",
    "WHITESPACE", // "excluding" '\"n"'
    "IO_NUMBER",

    "MIN_PUNCTUATION",
    "PIPE",       // |
    "AMPERSAND",  // &
    "SEMI",       // ;
    "LESS",       // <
    "GREATER",    // >
    "LPAREN",     // (
    "RPAREN",     // )
    "DOLLAR",     // $
    "BACKTICK",   // `
    "BACKSLASH",  /* \ */
    "QUOTE",      // '
    "DQUOTE",     // "
    "PLUS",       // +
    "DASH",       // -
    "ASTERISK",   // *
    "SLASH",      // /
    "QUESTION",   // ?
    "LBRACKET",   // [
    "RBRACKET",   // ]
    "HASH",       // #
    "TILDE",      // ~
    "EQUALS",     // =
    "PERCENT",    // %
    "LBRACE",     // {
    "RBRACE",     // }
    "BANG",       // !

    "AND_IF",     // &&
    "OR_IF",      // ||
    "DSEMI",      // ;;
    "DLESS",      // <<
    "DGREAT",     // >>
    "LESSAND",    // <&
    "GREATAND",   // >&
    "LESSGREAT",  // <>
    "DLESSDASH",  // <<-
    "CLOBBER",    // >|
    "DLBRACKET",  // [[
    "DRBRACKET",  // ]]
    "MAX_PUNCTUATION",

    "MIN_RESERVED_WORD",
    "IF",
    "THEN",
    "ELSE",
    "ELIF",
    "FI",
    "DO",
    "DONE",
    "CASE",
    "ESAC",
    "WHILE",
    "UNTIL",
    "FOR",
    "IN",
    "FUNCTION",
    "SELECT",
    "MAX_RESERVED_WORD"
};

char* tokentype_to_str(TokenType type) {
    int n = (int) type;
    assert(0 <= n && n <= N_TOKEN_TYPES);
    return TTYPE_TO_STR[n];
}
