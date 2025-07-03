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
	Selectors    []string
	SelectedData map[string]interface{} // Pour stocker les données sélectionnées
	NextURL      string
	Finished     bool
}

// PageInfo contient toutes les informations structurées de la page
type PageInfo struct {
	URL        string
	Title      string
	H1         []string
	H2         []string
	H3         []string
	Paragraphs []string
	Links      []parser.Link
	Images     []ImageInfo
	Lists      []string
}

// ImageInfo contient les informations d'une image
type ImageInfo struct {
	Src string
	Alt string
}

// ElementCategory représente une catégorie d'éléments sélectionnables
type ElementCategory struct {
	Name        string
	Icon        string
	Key         string
	Elements    []string
	Description string
}

// PromptSelectors enters an interactive session where the user can pick elements
// by category and see a structured preview of the page content.
func PromptSelectors(root *html.Node, currentURL *url.URL) (TuiResult, error) {
	pageInfo := extractPageInfo(root, currentURL)
	reader := bufio.NewReader(os.Stdin)

	for {
		printStructuredPage(pageInfo)
		printSelectionMenu()

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
			result, err := handleLinkNavigation(line, pageInfo.Links)
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				continue
			}
			return result, nil

		case strings.ToLower(line) == "all" || strings.ToLower(line) == "tout":
			return handleSelectAll(pageInfo), nil

		default:
			result, err := handleCategorySelection(line, pageInfo)
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

func extractPageInfo(root *html.Node, currentURL *url.URL) PageInfo {
	info := PageInfo{
		URL: currentURL.String(),
	}

	// Extraire le titre
	titleNodes := parser.FindAll(root, "title")
	if len(titleNodes) > 0 {
		info.Title = strings.TrimSpace(parser.TextContent(titleNodes[0]))
	}

	// Extraire les H1, H2, H3
	h1Nodes := parser.FindAll(root, "h1")
	for _, node := range h1Nodes {
		text := strings.TrimSpace(parser.TextContent(node))
		if text != "" {
			info.H1 = append(info.H1, text)
		}
	}

	h2Nodes := parser.FindAll(root, "h2")
	for _, node := range h2Nodes {
		text := strings.TrimSpace(parser.TextContent(node))
		if text != "" {
			info.H2 = append(info.H2, text)
		}
	}

	h3Nodes := parser.FindAll(root, "h3")
	for _, node := range h3Nodes {
		text := strings.TrimSpace(parser.TextContent(node))
		if text != "" {
			info.H3 = append(info.H3, text)
		}
	}

	// Extraire les paragraphes
	pNodes := parser.FindAll(root, "p")
	for _, node := range pNodes {
		text := strings.TrimSpace(parser.TextContent(node))
		if text != "" && len(text) > 10 { // Ignorer les paragraphes trop courts
			info.Paragraphs = append(info.Paragraphs, text)
		}
	}

	// Extraire les liens avec URLs absolues
	links := parser.FindLinks(root)
	for _, link := range links {
		if link.Href != "" {
			parsedHref, err := url.Parse(link.Href)
			if err == nil {
				link.Href = currentURL.ResolveReference(parsedHref).String()
			}
			if link.Text != "" {
				info.Links = append(info.Links, link)
			}
		}
	}

	// Extraire les images
	imgNodes := parser.FindAll(root, "img")
	for _, node := range imgNodes {
		var src, alt string
		for _, attr := range node.Attr {
			if attr.Key == "src" {
				src = attr.Val
			}
			if attr.Key == "alt" {
				alt = attr.Val
			}
		}
		if src != "" {
			info.Images = append(info.Images, ImageInfo{Src: src, Alt: alt})
		}
	}

	// Extraire les listes
	listNodes := append(parser.FindAll(root, "ul"), parser.FindAll(root, "ol")...)
	for _, listNode := range listNodes {
		liNodes := parser.FindAll(listNode, "li")
		var listItems []string
		for _, li := range liNodes {
			text := strings.TrimSpace(parser.TextContent(li))
			if text != "" {
				listItems = append(listItems, text)
			}
		}
		if len(listItems) > 0 {
			info.Lists = append(info.Lists, strings.Join(listItems, " | "))
		}
	}

	return info
}

func printStructuredPage(info PageInfo) {
	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("📄 Page: %s\n\n", info.URL)

	// Titre de la page
	if info.Title != "" {
		fmt.Printf("🌐 Title: %s\n", info.Title)
	}

	// H1
	if len(info.H1) > 0 {
		fmt.Printf("🔠 H1:\n")
		for _, h1 := range info.H1 {
			fmt.Printf(" - %s\n", h1)
		}
	}

	// H2
	if len(info.H2) > 0 {
		fmt.Printf("📰 H2:\n")
		for _, h2 := range info.H2 {
			fmt.Printf(" - %s\n", h2)
		}
	}

	// H3
	if len(info.H3) > 0 {
		fmt.Printf("📋 H3:\n")
		for _, h3 := range info.H3 {
			fmt.Printf(" - %s\n", h3)
		}
	}

	// Paragraphes
	if len(info.Paragraphs) > 0 {
		fmt.Printf("📝 Paragraphs:\n")
		for i, p := range info.Paragraphs {
			if i < 5 { // Limiter l'affichage à 5 paragraphes
				preview := p
				if len(preview) > 100 {
					preview = preview[:97] + "..."
				}
				fmt.Printf(" - %s\n", preview)
			}
		}
		if len(info.Paragraphs) > 5 {
			fmt.Printf(" ... et %d autres paragraphes\n", len(info.Paragraphs)-5)
		}
	}

	// Liens
	if len(info.Links) > 0 {
		fmt.Printf("🔗 Links:\n")
		for i, link := range info.Links {
			if i < 5 { // Limiter l'affichage à 5 liens
				fmt.Printf(" - %s (%s)\n", link.Text, link.Href)
			}
		}
		if len(info.Links) > 5 {
			fmt.Printf(" ... et %d autres liens\n", len(info.Links)-5)
		}
	}

	// Images
	if len(info.Images) > 0 {
		fmt.Printf("🖼️ Images:\n")
		for i, img := range info.Images {
			if i < 3 { // Limiter l'affichage à 3 images
				alt := img.Alt
				if alt == "" {
					alt = "Sans description"
				}
				fmt.Printf(" - %s (%s)\n", alt, img.Src)
			}
		}
		if len(info.Images) > 3 {
			fmt.Printf(" ... et %d autres images\n", len(info.Images)-3)
		}
	}

	// Listes
	if len(info.Lists) > 0 {
		fmt.Printf("📄 Lists:\n")
		for i, list := range info.Lists {
			if i < 3 { // Limiter l'affichage à 3 listes
				preview := list
				if len(preview) > 80 {
					preview = preview[:77] + "..."
				}
				fmt.Printf(" - %s\n", preview)
			}
		}
		if len(info.Lists) > 3 {
			fmt.Printf(" ... et %d autres listes\n", len(info.Lists)-3)
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
}

func printSelectionMenu() {
	fmt.Printf("\n📝 SÉLECTIONNER LES ÉLÉMENTS À EXTRAIRE :\n")
	fmt.Printf("  [title]     🌐 Titre de la page\n")
	fmt.Printf("  [h1]        🔠 Tous les titres H1\n")
	fmt.Printf("  [h2]        📰 Tous les titres H2\n")
	fmt.Printf("  [h3]        📋 Tous les titres H3\n")
	fmt.Printf("  [p]         📝 Tous les paragraphes\n")
	fmt.Printf("  [links]     🔗 Tous les liens\n")
	fmt.Printf("  [images]    🖼️ Toutes les images\n")
	fmt.Printf("  [lists]     📄 Toutes les listes\n")
	fmt.Printf("  [all]       ✨ Tous les éléments\n")
	fmt.Printf("  [L0,L1...]  🌐 Naviguer vers un lien (L0 = premier lien)\n")
	fmt.Printf("  [fini]      ✅ Terminer et générer le JSON\n")
	fmt.Printf("  [aide]      ❓ Afficher l'aide\n")
}

func printHelp() {
	fmt.Printf("\n" + strings.Repeat("*", 60) + "\n")
	fmt.Printf("🆘 AIDE DÉTAILLÉE\n")
	fmt.Printf(strings.Repeat("*", 60) + "\n")
	fmt.Printf("\n🎯 COMMENT UTILISER L'INTERFACE :\n")
	fmt.Printf("1. Examinez la structure de la page affichée ci-dessus\n")
	fmt.Printf("2. Sélectionnez les éléments que vous voulez extraire\n")
	fmt.Printf("3. Tapez 'fini' pour générer le JSON final\n")

	fmt.Printf("\n📋 EXEMPLES DE SÉLECTION :\n")
	fmt.Printf("  → 'title' pour extraire uniquement le titre\n")
	fmt.Printf("  → 'h1' pour extraire tous les H1\n")
	fmt.Printf("  → 'p' pour extraire tous les paragraphes\n")
	fmt.Printf("  → 'links' pour extraire tous les liens\n")
	fmt.Printf("  → 'all' pour extraire tous les éléments\n")

	fmt.Printf("\n🌐 NAVIGATION :\n")
	fmt.Printf("  → 'L0' pour aller au premier lien\n")
	fmt.Printf("  → 'L1' pour aller au deuxième lien, etc.\n")

	fmt.Printf("\n📤 RÉSULTAT :\n")
	fmt.Printf("Le JSON généré contiendra les clés correspondant aux éléments\n")
	fmt.Printf("sélectionnés (title, h1, paragraphs, links, etc.)\n")

	fmt.Printf("\n" + strings.Repeat("*", 60) + "\n")
}

func handleSelectAll(info PageInfo) TuiResult {
	selectedData := make(map[string]interface{})
	var selectors []string

	if info.Title != "" {
		selectedData["title"] = info.Title
		selectors = append(selectors, "title")
	}
	if len(info.H1) > 0 {
		selectedData["h1"] = info.H1
		selectors = append(selectors, "h1")
	}
	if len(info.H2) > 0 {
		selectedData["h2"] = info.H2
		selectors = append(selectors, "h2")
	}
	if len(info.H3) > 0 {
		selectedData["h3"] = info.H3
		selectors = append(selectors, "h3")
	}
	if len(info.Paragraphs) > 0 {
		selectedData["paragraphs"] = info.Paragraphs
		selectors = append(selectors, "p")
	}
	if len(info.Links) > 0 {
		linkUrls := make([]string, len(info.Links))
		for i, link := range info.Links {
			linkUrls[i] = link.Href
		}
		selectedData["links"] = linkUrls
		selectors = append(selectors, "links")
	}
	if len(info.Images) > 0 {
		imageSrcs := make([]string, len(info.Images))
		for i, img := range info.Images {
			imageSrcs[i] = img.Src
		}
		selectedData["images"] = imageSrcs
		selectors = append(selectors, "images")
	}
	if len(info.Lists) > 0 {
		selectedData["lists"] = info.Lists
		selectors = append(selectors, "lists")
	}

	fmt.Printf("✅ Tous les éléments sélectionnés !\n")
	return TuiResult{
		Selectors:    selectors,
		SelectedData: selectedData,
		Finished:     true,
	}
}

func handleCategorySelection(input string, info PageInfo) (*TuiResult, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	selectedData := make(map[string]interface{})
	var selectors []string

	switch input {
	case "title":
		if info.Title != "" {
			selectedData["title"] = info.Title
			selectors = append(selectors, "title")
			fmt.Printf("✅ Titre sélectionné: %s\n", info.Title)
		} else {
			return nil, fmt.Errorf("aucun titre trouvé sur cette page")
		}

	case "h1":
		if len(info.H1) > 0 {
			selectedData["h1"] = info.H1
			selectors = append(selectors, "h1")
			fmt.Printf("✅ %d titre(s) H1 sélectionné(s)\n", len(info.H1))
		} else {
			return nil, fmt.Errorf("aucun H1 trouvé sur cette page")
		}

	case "h2":
		if len(info.H2) > 0 {
			selectedData["h2"] = info.H2
			selectors = append(selectors, "h2")
			fmt.Printf("✅ %d titre(s) H2 sélectionné(s)\n", len(info.H2))
		} else {
			return nil, fmt.Errorf("aucun H2 trouvé sur cette page")
		}

	case "h3":
		if len(info.H3) > 0 {
			selectedData["h3"] = info.H3
			selectors = append(selectors, "h3")
			fmt.Printf("✅ %d titre(s) H3 sélectionné(s)\n", len(info.H3))
		} else {
			return nil, fmt.Errorf("aucun H3 trouvé sur cette page")
		}

	case "p", "paragraphs":
		if len(info.Paragraphs) > 0 {
			selectedData["paragraphs"] = info.Paragraphs
			selectors = append(selectors, "p")
			fmt.Printf("✅ %d paragraphe(s) sélectionné(s)\n", len(info.Paragraphs))
		} else {
			return nil, fmt.Errorf("aucun paragraphe trouvé sur cette page")
		}

	case "links":
		if len(info.Links) > 0 {
			linkUrls := make([]string, len(info.Links))
			for i, link := range info.Links {
				linkUrls[i] = link.Href
			}
			selectedData["links"] = linkUrls
			selectors = append(selectors, "links")
			fmt.Printf("✅ %d lien(s) sélectionné(s)\n", len(info.Links))
		} else {
			return nil, fmt.Errorf("aucun lien trouvé sur cette page")
		}

	case "images":
		if len(info.Images) > 0 {
			imageSrcs := make([]string, len(info.Images))
			for i, img := range info.Images {
				imageSrcs[i] = img.Src
			}
			selectedData["images"] = imageSrcs
			selectors = append(selectors, "images")
			fmt.Printf("✅ %d image(s) sélectionnée(s)\n", len(info.Images))
		} else {
			return nil, fmt.Errorf("aucune image trouvée sur cette page")
		}

	case "lists":
		if len(info.Lists) > 0 {
			selectedData["lists"] = info.Lists
			selectors = append(selectors, "lists")
			fmt.Printf("✅ %d liste(s) sélectionnée(s)\n", len(info.Lists))
		} else {
			return nil, fmt.Errorf("aucune liste trouvée sur cette page")
		}

	default:
		return nil, fmt.Errorf("sélection '%s' non reconnue. Tapez 'aide' pour voir les options", input)
	}

	return &TuiResult{
		Selectors:    selectors,
		SelectedData: selectedData,
		Finished:     true,
	}, nil
}

func handleFinish() (TuiResult, error) {
	fmt.Printf("✅ Session terminée.\n")
	return TuiResult{Finished: true}, nil
}

func handleLinkNavigation(input string, links []parser.Link) (TuiResult, error) {
	if len(input) < 2 {
		return TuiResult{}, fmt.Errorf("format invalide. Utilisez L suivi d'un numéro (ex: L0)")
	}

	idxStr := input[1:]
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return TuiResult{}, fmt.Errorf("numéro invalide: %s", idxStr)
	}

	if idx < 0 || idx >= len(links) {
		return TuiResult{}, fmt.Errorf("index %d hors limites (0-%d)", idx, len(links)-1)
	}

	fmt.Printf("🌐 Navigation vers: %s\n", links[idx].Href)
	return TuiResult{NextURL: links[idx].Href}, nil
}

