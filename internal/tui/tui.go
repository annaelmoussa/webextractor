package tui

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"webextractor/internal/parser"

	"golang.org/x/net/html"
)

// TuiResult holds the outcome of an interactive prompt session.
type TuiResult struct {
	Selectors []string
	NextURL   string
	Finished  bool
}

// ElementInfo contient les informations d'un élément avec sa catégorie
type ElementInfo struct {
	Node        *html.Node
	Category    string
	Description string
	Preview     string
	Selector    string
}

// PromptSelectors enters an interactive session where the user can pick nodes
// by numeric indices or navigate to a new page.
func PromptSelectors(root *html.Node, currentURL *url.URL) (TuiResult, error) {
	elements := categorizeElements(root)
	links := parser.FindLinks(root)

	// Make links absolute
	for i, link := range links {
		parsedHref, err := url.Parse(link.Href)
		if err == nil {
			links[i].Href = currentURL.ResolveReference(parsedHref).String()
		}
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		printHeader(currentURL.String())

		if len(elements) == 0 {
			fmt.Println("❌ Aucun élément extractible trouvé sur cette page.")
		} else {
			printElementsMenu(elements)
		}

		if len(links) > 0 {
			printLinksMenu(links)
		}

		printInstructions()

		fmt.Print("\n🎯 Votre choix : ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch {
		case strings.ToLower(line) == "aide" || strings.ToLower(line) == "help" || line == "?":
			printHelp()
			continue

		case strings.ToLower(line) == "fini" || strings.ToLower(line) == "done" || strings.ToLower(line) == "terminer":
			return handleFinish()

		case strings.HasPrefix(strings.ToLower(line), "l"):
			result, err := handleLinkNavigation(line, links)
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				continue
			}
			return result, nil

		case strings.ToLower(line) == "apercu" || strings.ToLower(line) == "preview":
			printPreview(elements)
			continue

		default:
			result, err := handleElementSelection(line, elements)
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				continue
			}
			if result != nil {
				return *result, nil
			}
		}
	}
}

func printHeader(url string) {
	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("🌐 EXTRACTION DE DONNÉES - %s\n", url)
	fmt.Printf(strings.Repeat("=", 70) + "\n")
}

func printElementsMenu(elements []ElementInfo) {
	categories := make(map[string][]ElementInfo)
	for _, elem := range elements {
		categories[elem.Category] = append(categories[elem.Category], elem)
	}

	fmt.Printf("\n📋 ÉLÉMENTS DISPONIBLES POUR EXTRACTION :\n")

	index := 0
	categoryIcons := map[string]string{
		"Titres":     "📰",
		"Textes":     "📄",
		"Liens":      "🔗",
		"Boutons":    "🔘",
		"Images":     "🖼️",
		"Conteneurs": "📦",
		"Autres":     "⚪",
	}

	categoryOrder := []string{"Titres", "Textes", "Liens", "Boutons", "Images", "Conteneurs", "Autres"}

	for _, category := range categoryOrder {
		if elements, ok := categories[category]; ok && len(elements) > 0 {
			icon, exists := categoryIcons[category]
			if !exists {
				icon = "⚪"
			}
			fmt.Printf("\n  %s %s :\n", icon, category)
			for _, elem := range elements {
				fmt.Printf("    [%d] %s\n", index, elem.Description)
				fmt.Printf("        💬 \"%s\"\n", elem.Preview)
				index++
			}
		}
	}
}

func printLinksMenu(links []parser.Link) {
	fmt.Printf("\n🔗 NAVIGATION VERS D'AUTRES PAGES :\n")
	for i, link := range links {
		linkText := link.Text
		if linkText == "" {
			linkText = "Lien sans texte"
		}
		if len(linkText) > 50 {
			linkText = linkText[:47] + "..."
		}
		fmt.Printf("  [L%d] %s\n", i, linkText)
		fmt.Printf("       🌐 %s\n", link.Href)
	}
}

func printInstructions() {
	fmt.Printf("\n" + strings.Repeat("-", 50) + "\n")
	fmt.Printf("📝 INSTRUCTIONS :\n")
	fmt.Printf("  • Tapez un numéro (ex: 0) pour sélectionner un élément\n")
	fmt.Printf("  • Tapez plusieurs numéros (ex: 0,2,5) pour sélectionner plusieurs éléments\n")
	fmt.Printf("  • Tapez une plage (ex: 0-3) pour sélectionner des éléments consécutifs\n")
	fmt.Printf("  • Tapez L suivi d'un numéro (ex: L0) pour naviguer vers un lien\n")
	fmt.Printf("  • Tapez 'apercu' pour voir ce qui serait extrait\n")
	fmt.Printf("  • Tapez 'aide' pour plus d'informations\n")
	fmt.Printf("  • Tapez 'fini' pour terminer la sélection\n")
	fmt.Printf(strings.Repeat("-", 50))
}

func printHelp() {
	fmt.Printf("\n" + strings.Repeat("*", 60) + "\n")
	fmt.Printf("🆘 AIDE DÉTAILLÉE\n")
	fmt.Printf(strings.Repeat("*", 60) + "\n")
	fmt.Printf("\n🎯 QU'EST-CE QUE L'EXTRACTION DE DONNÉES ?\n")
	fmt.Printf("WebExtractor vous aide à récupérer automatiquement du texte\n")
	fmt.Printf("ou des informations spécifiques d'une page web.\n")

	fmt.Printf("\n📋 TYPES D'ÉLÉMENTS :\n")
	fmt.Printf("  📰 Titres : Titres principaux et sous-titres de la page\n")
	fmt.Printf("  📄 Textes : Paragraphes et contenus textuels\n")
	fmt.Printf("  🔗 Liens : Liens vers d'autres pages ou ressources\n")
	fmt.Printf("  🔘 Boutons : Boutons cliquables et éléments interactifs\n")
	fmt.Printf("  🖼️ Images : Descriptions et textes alternatifs des images\n")
	fmt.Printf("  📦 Conteneurs : Sections qui regroupent d'autres éléments\n")

	fmt.Printf("\n💡 EXEMPLES D'UTILISATION :\n")
	fmt.Printf("  → Tapez '0' pour sélectionner le premier élément\n")
	fmt.Printf("  → Tapez '0,2,5' pour sélectionner les éléments 0, 2 et 5\n")
	fmt.Printf("  → Tapez '1-4' pour sélectionner les éléments 1, 2, 3 et 4\n")
	fmt.Printf("  → Tapez 'L0' pour aller vers le premier lien\n")

	fmt.Printf("\n✨ CONSEILS :\n")
	fmt.Printf("  • Commencez par sélectionner quelques éléments pour voir le résultat\n")
	fmt.Printf("  • Utilisez 'apercu' pour vérifier avant de terminer\n")
	fmt.Printf("  • Vous pouvez toujours naviguer vers d'autres pages et revenir\n")
	fmt.Printf("  • Les sélections sont cumulatives entre les pages\n")

	fmt.Printf("\n" + strings.Repeat("*", 60) + "\n")
}

func printPreview(elements []ElementInfo) {
	fmt.Printf("\n" + strings.Repeat("~", 50) + "\n")
	fmt.Printf("👀 APERÇU DE CE QUI SERAIT EXTRAIT :\n")
	fmt.Printf(strings.Repeat("~", 50) + "\n")

	if len(elements) == 0 {
		fmt.Printf("Aucun élément sélectionnable trouvé.\n")
		return
	}

	for i, elem := range elements {
		fmt.Printf("\n[%d] %s :\n", i, elem.Description)
		fmt.Printf("    Sélecteur CSS : %s\n", elem.Selector)
		fmt.Printf("    Contenu : \"%s\"\n", elem.Preview)
	}

	fmt.Printf("\n" + strings.Repeat("~", 50) + "\n")
}

func handleFinish() (TuiResult, error) {
	fmt.Printf("\n✅ Session terminée.\n")
	fmt.Printf("💡 Conseil : Si vous n'avez rien sélectionné, le programme va s'arrêter.\n")
	fmt.Printf("   Relancez avec des sélections pour extraire des données !\n")
	return TuiResult{Finished: true}, nil
}

func handleLinkNavigation(input string, links []parser.Link) (TuiResult, error) {
	idxStr := strings.TrimPrefix(strings.ToLower(input), "l")
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 0 || idx >= len(links) {
		return TuiResult{}, fmt.Errorf("Numéro de lien invalide. Utilisez L0, L1, etc. (disponibles : L0 à L%d)", len(links)-1)
	}

	fmt.Printf("\n🚀 Navigation vers : %s\n", links[idx].Href)
	return TuiResult{NextURL: links[idx].Href}, nil
}

func handleElementSelection(input string, elements []ElementInfo) (*TuiResult, error) {
	indices := parseIndices(input)
	if len(indices) == 0 {
		return nil, fmt.Errorf("Format invalide. Exemples valides : 0, 1,3,5, 0-2")
	}

	var selectors []string
	var selectedDescs []string

	for _, idx := range indices {
		if idx >= 0 && idx < len(elements) {
			selectors = append(selectors, elements[idx].Selector)
			selectedDescs = append(selectedDescs, elements[idx].Description)
		} else {
			return nil, fmt.Errorf("Index %d invalide. Éléments disponibles : 0 à %d", idx, len(elements)-1)
		}
	}

	if len(selectors) > 0 {
		fmt.Printf("\n✅ Sélectionné %d élément(s) :\n", len(selectors))
		for i, desc := range selectedDescs {
			fmt.Printf("  %d. %s\n", i+1, desc)
		}
		fmt.Printf("\n💾 Ces éléments seront extraits de cette page.\n")
		fmt.Printf("🔄 Vous pouvez continuer à naviguer ou taper 'fini' pour terminer.\n")

		return &TuiResult{Selectors: selectors}, nil
	}

	return nil, fmt.Errorf("Aucun élément valide sélectionné")
}

// findMeaningfulNodes collects element nodes that contain distinct text,
// filtering out parents that just wrap children with the same text content.
func findMeaningfulNodes(n *html.Node) []*html.Node {
	var nodes []*html.Node
	var rec func(*html.Node)

	rec = func(n *html.Node) {
		if n.Type != html.ElementNode {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
			return
		}

		// Skip non-user-friendly elements
		tag := strings.ToLower(n.Data)
		if isSkippableElement(tag) {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
			return
		}

		// Get text content of this node
		nodeText := strings.TrimSpace(previewText(n))

		// Skip empty nodes but continue recursion
		if nodeText == "" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
			return
		}

		// Check if this node has element children with the same text
		hasChildWithSameText := false
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				childText := strings.TrimSpace(previewText(c))
				if childText == nodeText {
					hasChildWithSameText = true
					break
				}
			}
		}

		if hasChildWithSameText {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
		} else {
			nodes = append(nodes, n)
		}
	}
	rec(n)
	return nodes
}

// isSkippableElement retourne true pour les éléments qu'un novice ne devrait pas voir
func isSkippableElement(tag string) bool {
	skippable := []string{
		"head", "meta", "title", "link", "style", "script", "noscript",
		"base", "object", "embed", "param", "source", "track", "area", "map",
		"colgroup", "col", "thead", "tbody", "tfoot", "option", "optgroup",
	}

	for _, skip := range skippable {
		if tag == skip {
			return true
		}
	}
	return false
}

func enumerateElements(n *html.Node, acc []*html.Node) []*html.Node {
	if n.Type == html.ElementNode {
		acc = append(acc, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		acc = enumerateElements(c, acc)
	}
	return acc
}

// previewText returns a short snippet of the node's text.
func previewText(n *html.Node) string {
	var b strings.Builder
	var rec func(*html.Node)
	rec = func(nd *html.Node) {
		if nd.Type == html.TextNode {
			trimmed := strings.TrimSpace(nd.Data)
			if trimmed != "" {
				b.WriteString(trimmed)
				b.WriteString(" ")
			}
		}
		for c := nd.FirstChild; c != nil; c = c.NextSibling {
			rec(c)
		}
	}
	rec(n)
	preview := strings.TrimSpace(b.String())
	if len(preview) > 40 {
		preview = preview[:40] + "..."
	}
	return preview
}

// parseIndices converts a string like "0,2,6-9" to a slice of ints.
func parseIndices(input string) []int {
	var out []int
	parts := strings.Split(input, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, "-") {
			rng := strings.SplitN(p, "-", 2)
			if len(rng) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(strings.TrimSpace(rng[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(rng[1]))
			if err1 != nil || err2 != nil {
				continue
			}
			if start > end {
				start, end = end, start
			}
			for i := start; i <= end; i++ {
				out = append(out, i)
			}
		} else {
			idx, err := strconv.Atoi(p)
			if err == nil {
				out = append(out, idx)
			}
		}
	}
	return out
}

// buildSelector creates a simple selector string prioritizing #id, then first class, otherwise tag.
func buildSelector(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val != "" {
			return "#" + a.Val
		}
	}
	for _, a := range n.Attr {
		if a.Key == "class" {
			classes := strings.Fields(a.Val)
			if len(classes) > 0 {
				return "." + classes[0]
			}
		}
	}
	return n.Data
}

// categorizeElements analyse les nœuds et les classe par catégorie
func categorizeElements(root *html.Node) []ElementInfo {
	nodes := findMeaningfulNodes(root)
	var elements []ElementInfo

	for _, node := range nodes {
		category, description := categorizeNode(node)
		preview := previewText(node)
		selector := buildSelector(node)

		if preview != "" {
			elements = append(elements, ElementInfo{
				Node:        node,
				Category:    category,
				Description: description,
				Preview:     preview,
				Selector:    selector,
			})
		}
	}

	return elements
}

// categorizeNode détermine la catégorie et la description d'un nœud
func categorizeNode(node *html.Node) (category, description string) {
	tag := strings.ToLower(node.Data)
	preview := previewText(node)

	switch tag {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		level := tag[1]
		return "Titres", fmt.Sprintf("Titre niveau %s", string(level))

	case "p":
		return "Textes", "Paragraphe de texte"

	case "a":
		return "Liens", "Lien hypertexte"

	case "button", "input":
		inputType := ""
		for _, attr := range node.Attr {
			if attr.Key == "type" {
				inputType = attr.Val
				break
			}
		}
		if inputType == "submit" || inputType == "button" || tag == "button" {
			return "Boutons", "Bouton cliquable"
		}
		return "Autres", fmt.Sprintf("Champ de saisie (%s)", inputType)

	case "img":
		return "Images", "Image (texte alternatif)"

	case "div", "section", "article", "main", "aside", "header", "footer", "nav":
		return "Conteneurs", fmt.Sprintf("Section (%s)", tag)

	case "span", "em", "strong", "b", "i":
		return "Textes", "Texte formaté"

	case "li":
		return "Textes", "Élément de liste"

	case "td", "th":
		return "Textes", "Cellule de tableau"

	default:
		// Si c'est un élément avec du texte significatif
		if len(strings.TrimSpace(preview)) > 10 {
			return "Textes", fmt.Sprintf("Contenu (%s)", tag)
		}
		return "Autres", fmt.Sprintf("Élément %s", tag)
	}
}
