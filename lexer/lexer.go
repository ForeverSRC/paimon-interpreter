package lexer

import (
	"github.com/ForeverSRC/paimon-interpreter/token"
	strutil "github.com/ForeverSRC/paimon-interpreter/util/str"
)

type Lexer struct {
	input string
	// position 输入字符串当前位置
	position int
	// readPosition 输入字符串的当前读取位置（当前字符后一个字符）
	readPosition int
	// ch 当前正在查看的字符
	// 如果想要支持unicode，此处应该为rune，并且需要修改readChar的逻辑
	ch byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	l.readChar()

	return l
}

// readChar 读取ascii 字符
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 0 ASCII code of NUL
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.NewWithString(token.EQ, literal)
		} else {
			tok = token.New(token.ASSIGN, l.ch)
		}
	case '+':
		tok = token.New(token.PLUS, l.ch)
	case '-':
		tok = token.New(token.MINUS, l.ch)
	case '*':
		tok = token.New(token.ASTERISK, l.ch)
	case '/':
		tok = token.New(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.NewWithString(token.NOTEQ, literal)
		} else {
			tok = token.New(token.BANG, l.ch)
		}
	case '<':
		tok = token.New(token.LT, l.ch)
	case '>':
		tok = token.New(token.GT, l.ch)
	case ';':
		tok = token.New(token.SEMICOLON, l.ch)
	case '(':
		tok = token.New(token.LPAREN, l.ch)
	case ')':
		tok = token.New(token.RPAREN, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case '{':
		tok = token.New(token.LBARCE, l.ch)
	case '}':
		tok = token.New(token.RBARCE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if strutil.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if strutil.IsDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for strutil.IsLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for strutil.IsDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
