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
	SelectedData map[string]interface{}
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

// SelectableElement représente un élément individuel sélectionnable
type SelectableElement struct {
	Index       int
	Type        string
	Content     string
	FullContent string
	Data        interface{}
}

// SelectionState garde l'état des sélections en cours
type SelectionState struct {
	Elements []SelectableElement
	Selected []bool
	PageInfo PageInfo
}

// PromptSelectors enters an interactive session where the user can pick elements
// individually and combine selections across different categories.
func PromptSelectors(root *html.Node, currentURL *url.URL) (TuiResult, error) {
	pageInfo := extractPageInfo(root, currentURL)
	elements := buildSelectableElements(pageInfo)
	state := SelectionState{
		Elements: elements,
		Selected: make([]bool, len(elements)),
		PageInfo: pageInfo,
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		printSelectableElements(state)
		printSelectionStatus(state)
		printSelectionMenu()

		fmt.Print("\n🎯 Votre choix : ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch {
		case strings.ToLower(line) == "aide" || strings.ToLower(line) == "help" || line == "?":
			printHelp()
			continue

		case strings.ToLower(line) == "fini" || strings.ToLower(line) == "done" || strings.ToLower(line) == "terminer":
			return handleFinishWithSelections(state)

		case strings.ToLower(line) == "reset" || strings.ToLower(line) == "clear":
			for i := range state.Selected {
				state.Selected[i] = false
			}
			fmt.Printf("🔄 Toutes les sélections effacées\n")
			continue

		case strings.HasPrefix(strings.ToLower(line), "l"):
			result, err := handleLinkNavigation(line, pageInfo.Links)
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				continue
			}
			return result, nil

		case strings.ToLower(line) == "all" || strings.ToLower(line) == "tout":
			for i := range state.Selected {
				state.Selected[i] = true
			}
			fmt.Printf("✅ Tous les éléments sélectionnés !\n")
			continue

		default:
			err := handleIndexSelection(line, &state)
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				continue
			}
		}
	}
}

func extractPageInfo(root *html.Node, currentURL *url.URL) PageInfo {
	info := PageInfo{
		URL: currentURL.String(),
	}

	titleNodes := parser.FindAll(root, "title")
	if len(titleNodes) > 0 {
		info.Title = strings.TrimSpace(parser.TextContent(titleNodes[0]))
	}

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

	pNodes := parser.FindAll(root, "p")
	for _, node := range pNodes {
		text := strings.TrimSpace(parser.TextContent(node))
		if text != "" && len(text) > 10 {
			info.Paragraphs = append(info.Paragraphs, text)
		}
	}

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

func buildSelectableElements(info PageInfo) []SelectableElement {
	var elements []SelectableElement
	index := 0

	if info.Title != "" {
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "title",
			Content:     truncateText(info.Title, 80),
			FullContent: info.Title,
			Data:        info.Title,
		})
		index++
	}

	for _, h1 := range info.H1 {
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "h1",
			Content:     truncateText(h1, 80),
			FullContent: h1,
			Data:        h1,
		})
		index++
	}

	for _, h2 := range info.H2 {
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "h2",
			Content:     truncateText(h2, 80),
			FullContent: h2,
			Data:        h2,
		})
		index++
	}

	for _, h3 := range info.H3 {
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "h3",
			Content:     truncateText(h3, 80),
			FullContent: h3,
			Data:        h3,
		})
		index++
	}

	for _, p := range info.Paragraphs {
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "p",
			Content:     truncateText(p, 80),
			FullContent: p,
			Data:        p,
		})
		index++
	}

	for _, link := range info.Links {
		linkText := fmt.Sprintf("%s (%s)", link.Text, link.Href)
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "link",
			Content:     truncateText(linkText, 80),
			FullContent: linkText,
			Data:        link,
		})
		index++
	}

	for _, img := range info.Images {
		alt := img.Alt
		if alt == "" {
			alt = "Sans description"
		}
		imgText := fmt.Sprintf("%s (%s)", alt, img.Src)
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "image",
			Content:     truncateText(imgText, 80),
			FullContent: imgText,
			Data:        img,
		})
		index++
	}

	for _, list := range info.Lists {
		elements = append(elements, SelectableElement{
			Index:       index,
			Type:        "list",
			Content:     truncateText(list, 80),
			FullContent: list,
			Data:        list,
		})
		index++
	}

	return elements
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func printSelectableElements(state SelectionState) {
	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("📄 Page: %s\n", state.PageInfo.URL)
	fmt.Printf("📝 ÉLÉMENTS SÉLECTIONNABLES :\n")
	fmt.Printf(strings.Repeat("=", 70) + "\n")

	for i, elem := range state.Elements {
		selectedMark := " "
		if state.Selected[i] {
			selectedMark = "✅"
		}

		var icon string
		switch elem.Type {
		case "title":
			icon = "🌐"
		case "h1":
			icon = "🔠"
		case "h2":
			icon = "📰"
		case "h3":
			icon = "📋"
		case "p":
			icon = "📝"
		case "link":
			icon = "🔗"
		case "image":
			icon = "🖼️"
		case "list":
			icon = "📄"
		default:
			icon = "📌"
		}

		fmt.Printf("%s [%2d] %s %s %s\n", selectedMark, i, icon, strings.ToUpper(elem.Type), elem.Content)
	}
	fmt.Printf(strings.Repeat("=", 70) + "\n")
}

func printSelectionStatus(state SelectionState) {
	selectedCount := 0
	for _, selected := range state.Selected {
		if selected {
			selectedCount++
		}
	}

	if selectedCount > 0 {
		fmt.Printf("\n📊 SÉLECTIONS ACTUELLES : %d élément(s) sélectionné(s)\n", selectedCount)

		typeCount := make(map[string]int)
		for i, elem := range state.Elements {
			if state.Selected[i] {
				typeCount[elem.Type]++
			}
		}

		var parts []string
		for elemType, count := range typeCount {
			var icon string
			switch elemType {
			case "title":
				icon = "🌐"
			case "h1":
				icon = "🔠"
			case "h2":
				icon = "📰"
			case "h3":
				icon = "📋"
			case "p":
				icon = "📝"
			case "link":
				icon = "🔗"
			case "image":
				icon = "🖼️"
			case "list":
				icon = "📄"
			}
			parts = append(parts, fmt.Sprintf("%s %s(%d)", icon, elemType, count))
		}
		fmt.Printf("  → %s\n", strings.Join(parts, ", "))
	} else {
		fmt.Printf("\n📊 Aucun élément sélectionné pour le moment\n")
	}
}

func printSelectionMenu() {
	fmt.Printf("\n📝 OPTIONS DE SÉLECTION :\n")
	fmt.Printf("  • Indices individuels : 0, 5, 12\n")
	fmt.Printf("  • Plages d'indices : 0-5, 10-15\n")
	fmt.Printf("  • Combinaisons : 0,3,7-9,15\n")
	fmt.Printf("  • [all]    ✨ Sélectionner tous les éléments\n")
	fmt.Printf("  • [reset]  🔄 Effacer toutes les sélections\n")
	fmt.Printf("  • [L0,L1]  🌐 Naviguer vers un lien\n")
	fmt.Printf("  • [fini]   ✅ Terminer et générer le JSON\n")
	fmt.Printf("  • [aide]   ❓ Afficher l'aide détaillée\n")
}

func printHelp() {
	fmt.Printf("\n" + strings.Repeat("*", 60) + "\n")
	fmt.Printf("🆘 AIDE DÉTAILLÉE - SÉLECTION GRANULAIRE\n")
	fmt.Printf(strings.Repeat("*", 60) + "\n")

	fmt.Printf("\n🎯 PRINCIPE :\n")
	fmt.Printf("Chaque élément de la page a un numéro [0, 1, 2, ...]\n")
	fmt.Printf("Vous pouvez sélectionner exactement les éléments que vous voulez.\n")

	fmt.Printf("\n📋 EXEMPLES DE SÉLECTION :\n")
	fmt.Printf("  → '0' pour sélectionner uniquement l'élément 0\n")
	fmt.Printf("  → '0,3,5' pour sélectionner les éléments 0, 3 et 5\n")
	fmt.Printf("  → '0-5' pour sélectionner les éléments 0 à 5 inclus\n")
	fmt.Printf("  → '0,3-7,10' pour combiner individuels et plages\n")
	fmt.Printf("  → 'all' pour sélectionner tous les éléments\n")
	fmt.Printf("  → 'reset' pour effacer toutes les sélections\n")

	fmt.Printf("\n💡 STRATÉGIE RECOMMANDÉE :\n")
	fmt.Printf("1. Examinez la liste numérotée des éléments\n")
	fmt.Printf("2. Notez les numéros des éléments qui vous intéressent\n")
	fmt.Printf("3. Sélectionnez-les par indices ou plages\n")
	fmt.Printf("4. Vérifiez vos sélections dans le résumé\n")
	fmt.Printf("5. Ajoutez/retirez des éléments si nécessaire\n")
	fmt.Printf("6. Tapez 'fini' pour générer le JSON\n")

	fmt.Printf("\n📤 RÉSULTAT :\n")
	fmt.Printf("Le JSON contiendra uniquement les éléments que vous avez\n")
	fmt.Printf("spécifiquement sélectionnés, organisés par type.\n")

	fmt.Printf("\n" + strings.Repeat("*", 60) + "\n")
}

func handleIndexSelection(input string, state *SelectionState) error {
	indices, err := parseIndices(input, len(state.Elements))
	if err != nil {
		return fmt.Errorf("format invalide: %v\nUtilisez des indices (0,1,2) ou plages (0-5)", err)
	}

	selectedCount := 0
	for _, idx := range indices {
		if !state.Selected[idx] {
			state.Selected[idx] = true
			selectedCount++
		}
	}

	fmt.Printf("✅ %d élément(s) sélectionné(s) : %v\n", selectedCount, indices)
	return nil
}

func parseIndices(input string, maxIndex int) ([]int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("entrée vide")
	}

	var indices []int
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("format de plage invalide: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("début de plage invalide: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("fin de plage invalide: %s", rangeParts[1])
			}

			if start < 0 || end >= maxIndex || start > end {
				return nil, fmt.Errorf("plage invalide: %d-%d (doit être entre 0 et %d)", start, end, maxIndex-1)
			}

			for i := start; i <= end; i++ {
				indices = append(indices, i)
			}
		} else {
			idx, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("index invalide: %s", part)
			}

			if idx < 0 || idx >= maxIndex {
				return nil, fmt.Errorf("index %d hors limites (0-%d)", idx, maxIndex-1)
			}

			indices = append(indices, idx)
		}
	}

	uniqueIndices := make([]int, 0, len(indices))
	seen := make(map[int]bool)
	for _, idx := range indices {
		if !seen[idx] {
			uniqueIndices = append(uniqueIndices, idx)
			seen[idx] = true
		}
	}

	return uniqueIndices, nil
}

func handleFinishWithSelections(state SelectionState) (TuiResult, error) {
	selectedData := make(map[string]interface{})
	var selectors []string

	selectedByType := make(map[string][]interface{})

	for i, elem := range state.Elements {
		if state.Selected[i] {
			switch elem.Type {
			case "title":
				selectedData["title"] = elem.Data.(string)
				selectors = append(selectors, "title")

			case "h1":
				selectedByType["h1"] = append(selectedByType["h1"], elem.Data)

			case "h2":
				selectedByType["h2"] = append(selectedByType["h2"], elem.Data)

			case "h3":
				selectedByType["h3"] = append(selectedByType["h3"], elem.Data)

			case "p":
				selectedByType["paragraphs"] = append(selectedByType["paragraphs"], elem.Data)

			case "link":
				link := elem.Data.(parser.Link)
				selectedByType["links"] = append(selectedByType["links"], link.Href)

			case "image":
				img := elem.Data.(ImageInfo)
				selectedByType["images"] = append(selectedByType["images"], img.Src)

			case "list":
				selectedByType["lists"] = append(selectedByType["lists"], elem.Data)
			}
		}
	}

	for elemType, items := range selectedByType {
		var stringSlice []string
		for _, item := range items {
			stringSlice = append(stringSlice, item.(string))
		}
		selectedData[elemType] = stringSlice
		selectors = append(selectors, elemType)
	}

	selectedCount := len(selectors)
	if _, hasTitle := selectedData["title"]; hasTitle {
		selectedCount += len(selectedByType) - 1
	} else {
		selectedCount = len(selectedByType)
	}

	if selectedCount == 0 {
		fmt.Printf("⚠️ Aucun élément sélectionné.\n")
		return TuiResult{Finished: false}, nil
	}

	fmt.Printf("✅ Session terminée avec %d élément(s) sélectionné(s)!\n", selectedCount)

	return TuiResult{
		Selectors:    selectors,
		SelectedData: selectedData,
		Finished:     true,
	}, nil
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
