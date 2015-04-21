import random
import string
import subprocess
import sys
import unittest

TEST_SHIM_PROG = None

class TestLex(unittest.TestCase):

    def run_lexer(self, text):
        return subprocess.check_output([TEST_SHIM_PROG, text])

    def test_lex_name(self):
        """A NAME is a WORD that does not start with a digit"""
        self.assertEqual(
            self.run_lexer('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_'),
            '("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_" NAME)')

    def test_lex_big_name(self):
        # put a letter at the front to ensure it's a NAME and not a WORD
        name = "a" + "".join([random.choice(string.letters + string.digits + '_') for _ in xrange(100000)])
        self.assertEqual(self.run_lexer(name), '("%s" NAME)' % name)

    def test_lex_word(self):
        """A WORD consists only of underscores, digits, and letters"""
        self.assertEqual(
            self.run_lexer('1234abc___'),
            '("1234abc___" WORD)')

    def test_lex_whitespace(self):
        self.assertEqual(
            self.run_lexer(' \t'),
            '(" \t" WHITESPACE)')

    def test_lex_newline(self):
        self.assertEqual(
            self.run_lexer('   \n   '),
            '("   " WHITESPACE)("\n" NEWLINE)("   " WHITESPACE)')

    def test_lex_words_names(self):
        self.assertEqual(
            self.run_lexer('abc 123 _123'),
            '("abc" NAME)'
            '(" " WHITESPACE)'
            '("123" WORD)'
            '(" " WHITESPACE)'
            '("_123" NAME)')

    def test_lex_greater_thans(self):
        self.assertEqual(
            self.run_lexer('< << <<-'),
            '("<" LESS)'
            '(" " WHITESPACE)'
            '("<<" DLESS)'
            '(" " WHITESPACE)'
            '("<<-" DLESSDASH)')

    def test_lex_greater_thans_no_whitespace(self):
        self.assertEqual(
            self.run_lexer('<<<<<-'),
            '("<<" DLESS)("<<" DLESS)("<" LESS)("-" DASH)')

    def test_lex_multiletter_ops(self):
        self.assertEqual(
            self.run_lexer('&& || ;; << >> <& >& <> <<- >| [[ ]]'),
                '("&&" AND_IF)(" " WHITESPACE)'
                '("||" OR_IF)(" " WHITESPACE)'
                '(";;" DSEMI)(" " WHITESPACE)'
                '("<<" DLESS)(" " WHITESPACE)'
                '(">>" DGREAT)(" " WHITESPACE)'
                '("<&" LESSAND)(" " WHITESPACE)'
                '(">&" GREATAND)(" " WHITESPACE)'
                '("<>" LESSGREAT)(" " WHITESPACE)'
                '("<<-" DLESSDASH)(" " WHITESPACE)'
                '(">|" CLOBBER)(" " WHITESPACE)'
                '("[[" DLBRACKET)(" " WHITESPACE)'
                '("]]" DRBRACKET)')

    def test_lex_single_letter_ops(self):
        self.assertEqual(
            self.run_lexer('| & ; < > ( ) $ ` \\ \' " + - * / ? [ ] # ~ = % { } !'),
                '("|" PIPE)(" " WHITESPACE)'
                '("&" AMPERSAND)(" " WHITESPACE)'
                '(";" SEMI)(" " WHITESPACE)'
                '("<" LESS)(" " WHITESPACE)'
                '(">" GREATER)(" " WHITESPACE)'
                '("(" LPAREN)(" " WHITESPACE)'
                '(")" RPAREN)(" " WHITESPACE)'
                '("$" DOLLAR)(" " WHITESPACE)'
                '("`" BACKTICK)(" " WHITESPACE)'
                '("\\" BACKSLASH)(" " WHITESPACE)'
                '("\'" QUOTE)(" " WHITESPACE)'
                '(""" DQUOTE)(" " WHITESPACE)'
                '("+" PLUS)(" " WHITESPACE)'
                '("-" DASH)(" " WHITESPACE)'
                '("*" ASTERISK)(" " WHITESPACE)'
                '("/" SLASH)(" " WHITESPACE)'
                '("?" QUESTION)(" " WHITESPACE)'
                '("[" LBRACKET)(" " WHITESPACE)'
                '("]" RBRACKET)(" " WHITESPACE)'
                '("#" HASH)(" " WHITESPACE)'
                '("~" TILDE)(" " WHITESPACE)'
                '("=" EQUALS)(" " WHITESPACE)'
                '("%" PERCENT)(" " WHITESPACE)'
                '("{" LBRACE)(" " WHITESPACE)'
                '("}" RBRACE)(" " WHITESPACE)'
                '("!" BANG)')

    def test_lex_reserved_words(self):
        self.assertEqual(
            self.run_lexer('if then else elif fi do done case esac while until '
                           'for in function select'),
            '("if" IF)(" " WHITESPACE)'
            '("then" THEN)(" " WHITESPACE)'
            '("else" ELSE)(" " WHITESPACE)'
            '("elif" ELIF)(" " WHITESPACE)'
            '("fi" FI)(" " WHITESPACE)'
            '("do" DO)(" " WHITESPACE)'
            '("done" DONE)(" " WHITESPACE)'
            '("case" CASE)(" " WHITESPACE)'
            '("esac" ESAC)(" " WHITESPACE)'
            '("while" WHILE)(" " WHITESPACE)'
            '("until" UNTIL)(" " WHITESPACE)'
            '("for" FOR)(" " WHITESPACE)'
            '("in" IN)(" " WHITESPACE)'
            '("function" FUNCTION)(" " WHITESPACE)'
            '("select" SELECT)')


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print "Usage: %s <test-shim-executable>" % sys.argv[0]
        sys.exit(1)
    TEST_SHIM_PROG = sys.argv[1]
    print "Using test shim %s" % TEST_SHIM_PROG
    sys.argv.pop(1)
    unittest.main()
