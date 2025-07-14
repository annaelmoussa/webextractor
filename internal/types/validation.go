package types

import (
	"fmt"
	"path/filepath"
	"strings"
)

// URLString représente une URL sous forme de chaîne validée
type URLString string

// NewURLString crée une nouvelle URL string après validation basique
func NewURLString(url string) (URLString, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Validation basique du format
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", fmt.Errorf("URL must start with http:// or https://")
	}

	return URLString(url), nil
}

// String retourne la représentation string
func (u URLString) String() string {
	return string(u)
}

// Scheme retourne le schéma de l'URL
func (u URLString) Scheme() string {
	if strings.HasPrefix(string(u), "https://") {
		return "https"
	}
	return "http"
}

// IsSecure retourne true si l'URL utilise HTTPS
func (u URLString) IsSecure() bool {
	return u.Scheme() == "https"
}

// CSSSelector représente un sélecteur CSS validé
type CSSSelector string

// NewCSSSelector crée un nouveau sélecteur CSS après validation
func NewCSSSelector(selector string) (CSSSelector, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return "", fmt.Errorf("CSS selector cannot be empty")
	}

	// Validation basique - peut être étendue
	if strings.Contains(selector, "<") || strings.Contains(selector, ">") {
		return "", fmt.Errorf("invalid characters in CSS selector")
	}

	return CSSSelector(selector), nil
}

// String retourne la représentation string
func (cs CSSSelector) String() string {
	return string(cs)
}

// IsClass retourne true si c'est un sélecteur de classe
func (cs CSSSelector) IsClass() bool {
	return strings.HasPrefix(string(cs), ".")
}

// IsID retourne true si c'est un sélecteur d'ID
func (cs CSSSelector) IsID() bool {
	return strings.HasPrefix(string(cs), "#")
}

// IsTag retourne true si c'est un sélecteur de tag
func (cs CSSSelector) IsTag() bool {
	return !cs.IsClass() && !cs.IsID()
}

// FilePath représente un chemin de fichier sécurisé
type FilePath string

// NewFilePath crée un nouveau chemin de fichier après validation
func NewFilePath(path string) (FilePath, error) {
	if path == "" || path == "-" {
		return FilePath(path), nil // stdout est autorisé
	}

	cleanPath := filepath.Clean(path)

	// Vérification de la traversée de répertoires
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("directory traversal not allowed in path: %s", path)
	}

	// Vérification des répertoires système
	if strings.HasPrefix(cleanPath, "/etc/") ||
		strings.HasPrefix(cleanPath, "/usr/") ||
		(strings.HasPrefix(cleanPath, "/var/") && !strings.HasPrefix(cleanPath, "/var/folders/")) {
		return "", fmt.Errorf("cannot write to system directories: %s", path)
	}

	return FilePath(cleanPath), nil
}

// String retourne la représentation string
func (fp FilePath) String() string {
	return string(fp)
}

// IsStdout retourne true si c'est stdout
func (fp FilePath) IsStdout() bool {
	return string(fp) == "-" || string(fp) == ""
}

// Extension retourne l'extension du fichier
func (fp FilePath) Extension() string {
	if fp.IsStdout() {
		return ""
	}
	return filepath.Ext(string(fp))
}

// HTTPStatus représente un code de statut HTTP
type HTTPStatus int

const (
	StatusOK            HTTPStatus = 200
	StatusNotFound      HTTPStatus = 404
	StatusInternalError HTTPStatus = 500
)

// IsSuccess retourne true si le statut indique un succès
func (hs HTTPStatus) IsSuccess() bool {
	return int(hs) >= 200 && int(hs) < 300
}

// IsClientError retourne true si c'est une erreur client
func (hs HTTPStatus) IsClientError() bool {
	return int(hs) >= 400 && int(hs) < 500
}

// IsServerError retourne true si c'est une erreur serveur
func (hs HTTPStatus) IsServerError() bool {
	return int(hs) >= 500 && int(hs) < 600
}

// String retourne la représentation string
func (hs HTTPStatus) String() string {
	return fmt.Sprintf("%d", int(hs))
}
