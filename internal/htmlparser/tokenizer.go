package htmlparser

import (
	"io"
	"strings"
	"unicode"
)

// TokenType represents the type of a token.
type TokenType int

const (
	ErrorToken TokenType = iota
	TextToken
	StartTagToken
	EndTagToken
	SelfClosingTagToken
	CommentToken
	DoctypeToken
)

// Token represents a token found during parsing.
type Token struct {
	Type TokenType
	Data string
	Attr []Attribute
}

// Tokenizer breaks HTML into tokens.
type Tokenizer struct {
	r     io.Reader
	raw   []byte
	pos   int
	token Token
}

// NewTokenizer creates a new tokenizer.
func NewTokenizer(r io.Reader) *Tokenizer {
	raw, _ := io.ReadAll(r)
	return &Tokenizer{
		r:   r,
		raw: raw,
		pos: 0,
	}
}

// Next advances to the next token.
func (t *Tokenizer) Next() TokenType {
	if t.pos >= len(t.raw) {
		return ErrorToken
	}
	
	if t.raw[t.pos] == '<' {
		return t.readTag()
	}
	return t.readText()
}

// Token returns the current token.
func (t *Tokenizer) Token() Token {
	return t.token
}


func (t *Tokenizer) readText() TokenType {
	start := t.pos
	for t.pos < len(t.raw) && t.raw[t.pos] != '<' {
		t.pos++
	}
	text := string(t.raw[start:t.pos])
	text = strings.TrimSpace(text)
	if text == "" {
		// Skip empty text and try next token
		if t.pos < len(t.raw) {
			return t.Next()
		}
		return ErrorToken
	}
	t.token = Token{
		Type: TextToken,
		Data: text,
	}
	return TextToken
}

func (t *Tokenizer) readTag() TokenType {
	if t.pos >= len(t.raw) || t.raw[t.pos] != '<' {
		return ErrorToken
	}
	
	t.pos++ // skip '<'
	
	if t.pos < len(t.raw) && t.raw[t.pos] == '!' {
		return t.readComment()
	}
	
	if t.pos < len(t.raw) && t.raw[t.pos] == '/' {
		return t.readEndTag()
	}
	
	return t.readStartTag()
}

func (t *Tokenizer) readComment() TokenType {
	if t.pos+3 >= len(t.raw) || string(t.raw[t.pos:t.pos+3]) != "!--" {
		// Handle <!DOCTYPE> and other declarations
		if t.pos+1 < len(t.raw) && strings.ToUpper(string(t.raw[t.pos:t.pos+2])) == "!D" {
			t.skipToEnd()
			t.token = Token{
				Type: DoctypeToken,
				Data: "doctype",
			}
			return DoctypeToken
		}
		t.skipToEnd()
		return ErrorToken
	}
	
	t.pos += 3 // skip "!--"
	start := t.pos
	
	for t.pos+2 < len(t.raw) {
		if string(t.raw[t.pos:t.pos+3]) == "-->" {
			comment := string(t.raw[start:t.pos])
			t.pos += 3
			t.token = Token{
				Type: CommentToken,
				Data: comment,
			}
			return CommentToken
		}
		t.pos++
	}
	
	t.skipToEnd()
	return ErrorToken
}

func (t *Tokenizer) readEndTag() TokenType {
	t.pos++ // skip '/'
	start := t.pos
	
	for t.pos < len(t.raw) && t.raw[t.pos] != '>' && !unicode.IsSpace(rune(t.raw[t.pos])) {
		t.pos++
	}
	
	tagName := strings.ToLower(string(t.raw[start:t.pos]))
	t.skipToEnd()
	
	t.token = Token{
		Type: EndTagToken,
		Data: tagName,
	}
	return EndTagToken
}

func (t *Tokenizer) readStartTag() TokenType {
	start := t.pos
	
	for t.pos < len(t.raw) && t.raw[t.pos] != '>' && t.raw[t.pos] != '/' && !unicode.IsSpace(rune(t.raw[t.pos])) {
		t.pos++
	}
	
	tagName := strings.ToLower(string(t.raw[start:t.pos]))
	
	t.skipWhitespace()
	
	var attrs []Attribute
	for t.pos < len(t.raw) && t.raw[t.pos] != '>' && t.raw[t.pos] != '/' {
		attr := t.readAttribute()
		if attr.Key != "" {
			attrs = append(attrs, attr)
		}
		t.skipWhitespace()
	}
	
	tokenType := StartTagToken
	if t.pos < len(t.raw) && t.raw[t.pos] == '/' {
		tokenType = SelfClosingTagToken
		t.pos++
	}
	
	t.skipToEnd()
	
	t.token = Token{
		Type: tokenType,
		Data: tagName,
		Attr: attrs,
	}
	return tokenType
}

func (t *Tokenizer) readAttribute() Attribute {
	t.skipWhitespace()
	
	if t.pos >= len(t.raw) || t.raw[t.pos] == '>' || t.raw[t.pos] == '/' {
		return Attribute{}
	}
	
	start := t.pos
	
	for t.pos < len(t.raw) && t.raw[t.pos] != '=' && t.raw[t.pos] != '>' && t.raw[t.pos] != '/' && !unicode.IsSpace(rune(t.raw[t.pos])) {
		t.pos++
	}
	
	key := strings.ToLower(string(t.raw[start:t.pos]))
	if key == "" {
		return Attribute{}
	}
	
	t.skipWhitespace()
	
	if t.pos >= len(t.raw) || t.raw[t.pos] != '=' {
		return Attribute{Key: key, Val: key}
	}
	
	t.pos++ // skip '='
	t.skipWhitespace()
	
	var val string
	if t.pos < len(t.raw) && (t.raw[t.pos] == '"' || t.raw[t.pos] == '\'') {
		quote := t.raw[t.pos]
		t.pos++
		start := t.pos
		for t.pos < len(t.raw) && t.raw[t.pos] != quote {
			t.pos++
		}
		val = string(t.raw[start:t.pos])
		if t.pos < len(t.raw) {
			t.pos++ // skip closing quote
		}
	} else {
		start := t.pos
		for t.pos < len(t.raw) && t.raw[t.pos] != '>' && t.raw[t.pos] != '/' && !unicode.IsSpace(rune(t.raw[t.pos])) {
			t.pos++
		}
		val = string(t.raw[start:t.pos])
	}
	
	return Attribute{Key: key, Val: val}
}

func (t *Tokenizer) skipWhitespace() {
	for t.pos < len(t.raw) && unicode.IsSpace(rune(t.raw[t.pos])) {
		t.pos++
	}
}

func (t *Tokenizer) skipToEnd() {
	for t.pos < len(t.raw) && t.raw[t.pos] != '>' {
		t.pos++
	}
	if t.pos < len(t.raw) {
		t.pos++ // skip '>'
	}
}