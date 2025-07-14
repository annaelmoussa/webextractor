package app

import (
	"fmt"
	"strings"

	"webextractor/internal/fetcher"
	"webextractor/internal/htmlparser"
	"webextractor/internal/io"
	"webextractor/internal/neturl"
	"webextractor/internal/parser"
	"webextractor/internal/tui"
	"webextractor/internal/types"
)

// App représente l'application WebExtractor
type App struct {
	config  *types.ExtractionConfig
	fetcher *fetcher.Fetcher
}

// New crée une nouvelle instance de l'application
func New(config *types.ExtractionConfig) *App {
	return &App{
		config:  config,
		fetcher: fetcher.New(config.Timeout),
	}
}

// Run exécute l'application
func (app *App) Run() error {
	if app.config.Selectors.IsEmpty() {
		if err := app.runInteractiveMode(); err != nil {
			return fmt.Errorf("interactive mode failed: %w", err)
		}
	}

	if app.config.Selectors.IsEmpty() && !app.config.IsStructuredMode() {
		printNoSelectionMessage()
		return nil
	}

	if app.config.IsStructuredMode() {
		return app.processStructuredOutput()
	}

	return app.processSelectorOutput()
}

// runInteractiveMode gère le mode interactif de sélection
func (app *App) runInteractiveMode() error {
	session := types.NewSessionState(app.config.URL)

	for {
		parsedURL, err := neturl.Parse(session.CurrentURL)
		if err != nil {
			return fmt.Errorf("invalid URL '%s': %w", session.CurrentURL, err)
		}

		fmt.Printf("Fetching %s...\n", session.CurrentURL)
		doc, err := app.fetcher.Fetch(session.CurrentURL)
		if err != nil {
			return fmt.Errorf("fetch error for %s: %w", session.CurrentURL, err)
		}

		res, err := tui.PromptSelectors(doc, parsedURL)
		if err != nil {
			return fmt.Errorf("TUI prompt failed: %w", err)
		}

		if len(res.Selectors) > 0 {
			if res.SelectedData != nil {
				session.SetStructuredData(res.SelectedData)
				app.config.SetStructuredMode(res.SelectedData)
				fmt.Printf("✅ Données structurées sélectionnées\n")
			} else {
				session.AddSelectors(res.Selectors)
				fmt.Printf("✅ Added %d selectors. Total: %d\n", len(res.Selectors), len(session.CollectedSelectors))
			}
		}

		if res.Finished {
			break
		}

		if res.NextURL != "" {
			session.CurrentURL = res.NextURL
		}
	}

	if !session.UseStructured {
		app.config.Selectors = session.FinalSelectors()
	}

	return nil
}

// processStructuredOutput traite la sortie en mode structuré
func (app *App) processStructuredOutput() error {
	structuredResult := convertToStructuredResult(app.config.URL, app.config.StructuredData)
	fmt.Printf("✅ Extraction terminée avec format structuré\n")
	printResultLocation(app.config.OutputPath)

	if err := io.WriteStructured(app.config.OutputPath.String(), structuredResult); err != nil {
		return fmt.Errorf("failed to write structured output: %w", err)
	}
	return nil
}

// processSelectorOutput traite la sortie en mode sélecteurs
func (app *App) processSelectorOutput() error {
	fmt.Printf("\n🔄 Extraction finale des données de %s...\n", app.config.URL)
	doc, err := app.fetcher.Fetch(app.config.URL)
	if err != nil {
		return fmt.Errorf("fetch error: %w", err)
	}

	results := extractUsingSelectors(doc, app.config.Selectors)
	extractionResult := types.NewExtractionResult(app.config.URL)
	extractionResult.SetMetrics(countTotalMatches(results), len(app.config.Selectors))

	fmt.Println(extractionResult.String())
	printResultLocation(app.config.OutputPath)

	if err := io.Write(app.config.OutputPath.String(), io.DocumentResult{URL: app.config.URL, Results: results}); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
}

// printNoSelectionMessage affiche un message quand aucune sélection n'est faite
func printNoSelectionMessage() {
	fmt.Println("\n🔄 Aucun élément sélectionné pour l'extraction.")
	fmt.Println("💡 Relancez le programme et sélectionnez des éléments pour extraire des données.")
	fmt.Println("📖 Exemple : go run main.go -url \"https://example.com\"")
}

// printResultLocation affiche où les résultats sont sauvegardés
func printResultLocation(outputPath types.OutputPath) {
	if outputPath.IsStdout() {
		fmt.Printf("📤 Résultats affichés ci-dessous :\n\n")
	} else {
		fmt.Printf("📁 Résultats sauvegardés dans : %s\n", outputPath.String())
	}
}

// convertToStructuredResult convertit les données brutes en résultat structuré
func convertToStructuredResult(url string, data map[string]any) io.StructuredResult {
	result := io.StructuredResult{URL: url}

	if title, ok := data["title"].(string); ok {
		result.Title = title
	}
	if h1List, ok := data["h1"].([]string); ok {
		result.H1 = h1List
	}
	if h2List, ok := data["h2"].([]string); ok {
		result.H2 = h2List
	}
	if h3List, ok := data["h3"].([]string); ok {
		result.H3 = h3List
	}
	if paragraphs, ok := data["paragraphs"].([]string); ok {
		result.Paragraphs = paragraphs
	}
	if links, ok := data["links"].([]string); ok {
		result.Links = links
	}
	if images, ok := data["images"].([]string); ok {
		result.Images = images
	}
	if lists, ok := data["lists"].([]string); ok {
		result.Lists = lists
	}

	return result
}

// extractUsingSelectors extrait les données en utilisant des sélecteurs CSS
func extractUsingSelectors(doc *htmlparser.Node, selectors types.SelectorList) []io.Result {
	var results []io.Result
	for _, sel := range selectors {
		sel = strings.TrimSpace(sel)
		if sel == "" {
			continue
		}
		nodes := parser.FindAll(doc, sel)
		var matches []string
		for _, n := range nodes {
			matches = append(matches, parser.TextContent(n))
		}
		results = append(results, io.Result{Selector: sel, Matches: matches})
	}
	return results
}

// countTotalMatches compte le nombre total de correspondances
func countTotalMatches(results []io.Result) int {
	total := 0
	for _, result := range results {
		total += len(result.Matches)
	}
	return total
}
