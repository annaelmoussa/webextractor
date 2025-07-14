package htmlparser

import (
	"io"
	"strings"
	"unicode"
)

// TokenType représente le type d'un token.
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

// Token représente un token trouvé pendant l'analyse.
type Token struct {
	Type TokenType   // On peut ajouter un type à un token exemple (TextToken, StartTagToken, EndTagToken, etc.)
	Data string      // On peut ajouter du texte à un token exemple (Hello World)
	Attr []Attribute // On peut ajouter des attributs à un token exemple (class, id, etc.)
}

// Tokenizer décompose le HTML en tokens.
type Tokenizer struct {
	raw   []byte // le HTML stocké en mémoire (tout le contenu)
	pos   int    // où on en est (curseur actuel dans raw)
	token Token  // le  dernier token trouvé
}

// NewTokenizer crée un nouveau tokenizer.
func NewTokenizer(r io.Reader) *Tokenizer {
	raw, _ := io.ReadAll(r)
	return &Tokenizer{
		raw: raw,
		pos: 0, // On commence à la position 0
	}
}

// Next avance au token suivant.
func (t *Tokenizer) Next() TokenType {
	// Si on est à la fin du HTML, on retourne une erreur
	if t.pos >= len(t.raw) {
		return ErrorToken
	}
	// Si le caractère courant est un <, alors on lit une balise
	if t.raw[t.pos] == '<' {
		// On lit une balise
		return t.readTag()
	}
	// Sinon on lit du texte
	return t.readText()
}

// Token retourne le token actuel.
func (t *Tokenizer) Token() Token {
	return t.token
}

func (t *Tokenizer) readText() TokenType {
	// On récupère la position de début du texte
	start := t.pos
	// On parcourt le HTML jusqu'à trouver un <
	for t.pos < len(t.raw) && t.raw[t.pos] != '<' {
		t.pos++
	}
	// On récupère le texte
	text := string(t.raw[start:t.pos])
	// On retire les espaces
	text = strings.TrimSpace(text)
	// Si le texte est vide, on ignore le texte vide et on essaie le token suivant
	if text == "" {
		// On ignore le texte vide et on essaie le token suivant
		if t.pos < len(t.raw) {
			return t.Next()
		}
		// Si on est à la fin du HTML, on retourne une erreur
		return ErrorToken
	}
	// On crée un nouveau token de type TextToken avec le texte trouvé
	t.token = Token{
		Type: TextToken,
		Data: text,
	}
	return TextToken
}

func (t *Tokenizer) readTag() TokenType {
	// Si on est à la fin du HTML ou si le caractère courant n'est pas un <, on retourne une erreur
	if t.pos >= len(t.raw) || t.raw[t.pos] != '<' {
		return ErrorToken
	}

	t.pos++ // On ignore le '<'
	// Si le caractère courant est un !, alors on lit une balise de commentaire
	if t.pos < len(t.raw) && t.raw[t.pos] == '!' {
		return t.readComment()
	}
	// Si le caractère courant est un /, alors on lit une balise de fermeture
	if t.pos < len(t.raw) && t.raw[t.pos] == '/' {
		return t.readEndTag()
	}

	return t.readStartTag()
}

func (t *Tokenizer) readComment() TokenType {
	// Si le caractère courant est un ! et que le suivant est un - et que le suivant est un -, alors on lit une balise de commentaire
	if t.pos+3 >= len(t.raw) || string(t.raw[t.pos:t.pos+3]) != "!--" {
		// On gère le <!DOCTYPE> et les autres déclarations
		// Si le caractère courant est un ! et que le suivant est un D, alors on lit une balise de type doctype
		if t.pos+1 < len(t.raw) && strings.ToUpper(string(t.raw[t.pos:t.pos+2])) == "!D" {
			t.skipToEnd()
			t.token = Token{
				Type: DoctypeToken,
				Data: "doctype",
			}
			return DoctypeToken
		}
		// On ignore le reste de la balise de commentaire
		t.skipToEnd()
		return ErrorToken
	}

	// On ignore le "!--"
	t.pos += 3
	// On récupère la position de début du commentaire
	start := t.pos
	// On parcourt le HTML jusqu'à trouver un -->
	for t.pos+2 < len(t.raw) {
		if string(t.raw[t.pos:t.pos+3]) == "-->" {
			// On récupère le commentaire
			comment := string(t.raw[start:t.pos])
			// On ignore le "-->"
			t.pos += 3
			// On crée un nouveau token de type CommentToken avec le commentaire trouvé
			t.token = Token{
				Type: CommentToken,
				Data: comment,
			}
			return CommentToken
		}
		// On avance d'un caractère
		t.pos++
	}

	// On ignore le reste de la balise de commentaire
	t.skipToEnd()
	return ErrorToken
}

func (t *Tokenizer) readEndTag() TokenType {
	// On ignore le '/'
	t.pos++
	// On récupère la position de début du tag
	start := t.pos
	// On parcourt le HTML jusqu'à trouver un >
	for t.pos < len(t.raw) && t.raw[t.pos] != '>' && !unicode.IsSpace(rune(t.raw[t.pos])) {
		t.pos++
	}

	// On récupère le nom du tag
	tagName := strings.ToLower(string(t.raw[start:t.pos]))
	// On ignore le reste de la balise de fermeture
	t.skipToEnd()
	// On crée un nouveau token de type EndTagToken avec le nom du tag trouvé
	t.token = Token{
		Type: EndTagToken,
		Data: tagName,
	}
	return EndTagToken
}

func (t *Tokenizer) readStartTag() TokenType {
	// On récupère la position de début du tag
	start := t.pos
	// On parcourt le HTML jusqu'à trouver un >
	for t.pos < len(t.raw) && t.raw[t.pos] != '>' && t.raw[t.pos] != '/' && !unicode.IsSpace(rune(t.raw[t.pos])) {
		t.pos++
	}
	// On récupère le nom du tag
	tagName := strings.ToLower(string(t.raw[start:t.pos]))
	// On ignore les espaces
	t.skipWhitespace()
	// On parcourt le HTML jusqu'à trouver un > ou un /
	var attrs []Attribute
	for t.pos < len(t.raw) && t.raw[t.pos] != '>' && t.raw[t.pos] != '/' {
		// On lit un attribut
		attr := t.readAttribute()
		if attr.Key != "" {
			attrs = append(attrs, attr)
		}
		t.skipWhitespace()
	}
	// On détermine le type de token
	tokenType := StartTagToken
	// Si le caractère courant est un /, alors on lit une balise auto-fermante
	if t.pos < len(t.raw) && t.raw[t.pos] == '/' {
		// On définit le type de token à SelfClosingTagToken
		tokenType = SelfClosingTagToken
		t.pos++
	}
	// On ignore le reste de la balise de début
	t.skipToEnd()
	// On crée un nouveau token de type StartTagToken avec le nom du tag, les attributs et le type de token
	t.token = Token{
		Type: tokenType,
		Data: tagName,
		Attr: attrs,
	}
	return tokenType
}

func (t *Tokenizer) readAttribute() Attribute {
	// On ignore les espaces
	t.skipWhitespace()
	// Si on est à la fin du HTML ou si le caractère courant est un > ou un /, alors on retourne un attribut vide
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

	t.pos++ // On ignore le '='
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
			t.pos++ // On ignore la quote de fermeture
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
		t.pos++ // On ignore le '>'
	}
}
