package types

import (
	"fmt"
	"strings"
	"time"

	"webextractor/internal/htmlparser"
	"webextractor/internal/neturl"
)

// SelectorList représente une liste de sélecteurs CSS
type SelectorList []string

// NewSelectorList crée une nouvelle liste de sélecteurs depuis une chaîne séparée par des virgules
func NewSelectorList(selectors string) SelectorList {
	if strings.TrimSpace(selectors) == "" {
		return SelectorList{}
	}
	
	var result SelectorList
	for _, sel := range strings.Split(selectors, ",") {
		if trimmed := strings.TrimSpace(sel); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// IsEmpty retourne true si la liste est vide
func (sl SelectorList) IsEmpty() bool {
	return len(sl) == 0
}

// Unique retourne une nouvelle liste sans doublons
func (sl SelectorList) Unique() SelectorList {
	seen := make(map[string]struct{})
	var result SelectorList
	
	for _, sel := range sl {
		if _, exists := seen[sel]; !exists {
			seen[sel] = struct{}{}
			result = append(result, sel)
		}
	}
	return result
}

// OutputPath représente un chemin de sortie sécurisé
type OutputPath string

const (
	StdoutPath OutputPath = "-"
)

// NewOutputPath crée un nouveau chemin de sortie
func NewOutputPath(path string) OutputPath {
	if path == "" {
		return StdoutPath
	}
	return OutputPath(path)
}

// IsStdout retourne true si la sortie est stdout
func (op OutputPath) IsStdout() bool {
	return op == StdoutPath
}

// String retourne la représentation string
func (op OutputPath) String() string {
	return string(op)
}

// ElementType représente le type d'un élément HTML
type ElementType string

const (
	ElementTypeTitle      ElementType = "title"
	ElementTypeH1         ElementType = "h1"
	ElementTypeH2         ElementType = "h2"
	ElementTypeH3         ElementType = "h3"
	ElementTypeParagraph  ElementType = "p"
	ElementTypeLink       ElementType = "link"
	ElementTypeImage      ElementType = "image"
	ElementTypeList       ElementType = "list"
)

// Icon retourne l'icône associée au type d'élément
func (et ElementType) Icon() string {
	switch et {
	case ElementTypeTitle:
		return "🌐"
	case ElementTypeH1:
		return "🔠"
	case ElementTypeH2:
		return "📰"
	case ElementTypeH3:
		return "📋"
	case ElementTypeParagraph:
		return "📝"
	case ElementTypeLink:
		return "🔗"
	case ElementTypeImage:
		return "🖼️"
	case ElementTypeList:
		return "📄"
	default:
		return "📌"
	}
}

// UserAgent représente un User-Agent HTTP
type UserAgent string

const (
	DefaultUserAgent UserAgent = "WebExtractor/0.1"
)

// String retourne la représentation string
func (ua UserAgent) String() string {
	return string(ua)
}

// ExtractionMode représente le mode d'extraction
type ExtractionMode int

const (
	ModeSelectorBased ExtractionMode = iota
	ModeStructured
)

// ExtractionConfig contient la configuration pour une extraction
type ExtractionConfig struct {
	URL         string
	Selectors   SelectorList
	OutputPath  OutputPath
	Timeout     time.Duration
	Mode        ExtractionMode
	StructuredData map[string]any
}

// NewExtractionConfig crée une nouvelle configuration d'extraction
func NewExtractionConfig(url string, selectors string, outputPath string, timeout time.Duration) *ExtractionConfig {
	return &ExtractionConfig{
		URL:        url,
		Selectors:  NewSelectorList(selectors),
		OutputPath: NewOutputPath(outputPath),
		Timeout:    timeout,
		Mode:       ModeSelectorBased,
	}
}

// SetStructuredMode configure le mode structuré
func (ec *ExtractionConfig) SetStructuredMode(data map[string]any) {
	ec.Mode = ModeStructured
	ec.StructuredData = data
}

// IsStructuredMode retourne true si le mode est structuré
func (ec *ExtractionConfig) IsStructuredMode() bool {
	return ec.Mode == ModeStructured
}

// SessionState représente l'état d'une session interactive
type SessionState struct {
	CurrentURL       string
	CollectedSelectors SelectorList
	StructuredData   map[string]any
	UseStructured    bool
}

// NewSessionState crée un nouvel état de session
func NewSessionState(startURL string) *SessionState {
	return &SessionState{
		CurrentURL:         startURL,
		CollectedSelectors: SelectorList{},
		UseStructured:      false,
	}
}

// AddSelectors ajoute des sélecteurs à la collection
func (ss *SessionState) AddSelectors(selectors []string) {
	ss.CollectedSelectors = append(ss.CollectedSelectors, selectors...)
}

// SetStructuredData configure les données structurées
func (ss *SessionState) SetStructuredData(data map[string]any) {
	ss.StructuredData = data
	ss.UseStructured = true
}

// FinalSelectors retourne les sélecteurs finaux uniques
func (ss *SessionState) FinalSelectors() SelectorList {
	return ss.CollectedSelectors.Unique()
}

// ExtractionResult représente le résultat d'une extraction
type ExtractionResult struct {
	URL          string
	TotalMatches int
	SelectorsUsed int
	Success      bool
	Error        error
}

// NewExtractionResult crée un nouveau résultat d'extraction
func NewExtractionResult(url string) *ExtractionResult {
	return &ExtractionResult{
		URL:     url,
		Success: true,
	}
}

// SetError marque le résultat comme échoué
func (er *ExtractionResult) SetError(err error) {
	er.Success = false
	er.Error = err
}

// SetMetrics configure les métriques du résultat
func (er *ExtractionResult) SetMetrics(totalMatches, selectorsUsed int) {
	er.TotalMatches = totalMatches
	er.SelectorsUsed = selectorsUsed
}

// String retourne une représentation string du résultat
func (er *ExtractionResult) String() string {
	if !er.Success {
		return fmt.Sprintf("❌ Extraction failed for %s: %v", er.URL, er.Error)
	}
	return fmt.Sprintf("✅ Extraction completed: %d selectors used, %d elements extracted", er.SelectorsUsed, er.TotalMatches)
}

// FetchRequest représente une requête de récupération
type FetchRequest struct {
	URL       string
	UserAgent UserAgent
	Timeout   time.Duration
}

// NewFetchRequest crée une nouvelle requête de récupération
func NewFetchRequest(url string, timeout time.Duration) *FetchRequest {
	return &FetchRequest{
		URL:       url,
		UserAgent: DefaultUserAgent,
		Timeout:   timeout,
	}
}

// FetchResult représente le résultat d'une récupération
type FetchResult struct {
	Document *htmlparser.Node
	URL      *neturl.URL
	Success  bool
	Error    error
}

// NewFetchResult crée un nouveau résultat de récupération
func NewFetchResult() *FetchResult {
	return &FetchResult{Success: true}
}

// SetError marque le résultat comme échoué
func (fr *FetchResult) SetError(err error) {
	fr.Success = false
	fr.Error = err
}

// SetDocument configure le document et l'URL
func (fr *FetchResult) SetDocument(doc *htmlparser.Node, url *neturl.URL) {
	fr.Document = doc
	fr.URL = url
}