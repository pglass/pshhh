#include "../pshscript/lex.h"
#include "../buf.h"
#include "../list.h"

int main(int argc, char** argv) {
    if (argc < 2) {
        printf("Usage: %s <text>\n", argv[0]);
        exit(1);
    }
    Buf* buf = buf_init(1024);
    buf_append(buf, argv[1], strlen(argv[1]));
    Lexer* lexer = lexer_init(buf);
    lexer_lex(lexer);
    ListNode* node = lexer->tokens->head->next;
    while (node != lexer->tokens->tail) {
        // token_print((Token*) node->item);
        Token* token = (Token*) node->item;
        printf("(\"%s\" %s)", token->data.str, tokentype_to_str(token->type));
        node = node->next;
    }
}
