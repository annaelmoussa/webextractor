package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"webextractor/internal/fetcher"
	"webextractor/internal/io"
	"webextractor/internal/parser"
	"webextractor/internal/tui"
)

func main() {
	urlPtr := flag.String("url", "", "URL of the web page to extract from (required)")
	selPtr := flag.String("sel", "", "CSS-like selector (tag, .class, #id). If omitted, interactive mode starts")
	outPtr := flag.String("out", "-", "Output JSON file path ('-' for stdout)")
	timeoutPtr := flag.Duration("timeout", 10*time.Second, "HTTP client timeout")

	flag.Parse()

	if *urlPtr == "" {
		fmt.Fprintln(os.Stderr, "❌ L'option -url est requise")
		flag.Usage()
		os.Exit(1)
	}

	f := fetcher.New(*timeoutPtr)

	var selectors []string
	var structuredData map[string]interface{}
	var useStructuredOutput bool

	if strings.TrimSpace(*selPtr) == "" {
		var err error
		selectors, structuredData, useStructuredOutput, err = interactiveSession(*urlPtr, f)
		if err != nil {
			log.Fatalf("❌ Erreur lors de la session interactive : %v", err)
		}
		if len(selectors) == 0 && !useStructuredOutput {
			fmt.Println("\n🔄 Aucun élément sélectionné pour l'extraction.")
			fmt.Println("💡 Relancez le programme et sélectionnez des éléments pour extraire des données.")
			fmt.Println("📖 Exemple : go run main.go -url \"https://example.com\"")
			os.Exit(0)
		}
	} else {
		selectors = strings.Split(*selPtr, ",")
		useStructuredOutput = false
	}

	if useStructuredOutput && structuredData != nil {
		structuredResult := io.StructuredResult{
			URL: *urlPtr,
		}

		if title, ok := structuredData["title"].(string); ok {
			structuredResult.Title = title
		}
		if h1List, ok := structuredData["h1"].([]string); ok {
			structuredResult.H1 = h1List
		}
		if h2List, ok := structuredData["h2"].([]string); ok {
			structuredResult.H2 = h2List
		}
		if h3List, ok := structuredData["h3"].([]string); ok {
			structuredResult.H3 = h3List
		}
		if paragraphs, ok := structuredData["paragraphs"].([]string); ok {
			structuredResult.Paragraphs = paragraphs
		}
		if links, ok := structuredData["links"].([]string); ok {
			structuredResult.Links = links
		}
		if images, ok := structuredData["images"].([]string); ok {
			structuredResult.Images = images
		}
		if lists, ok := structuredData["lists"].([]string); ok {
			structuredResult.Lists = lists
		}

		fmt.Printf("✅ Extraction terminée avec format structuré\n")

		if *outPtr == "-" {
			fmt.Printf("📤 Résultats affichés ci-dessous :\n\n")
		} else {
			fmt.Printf("📁 Résultats sauvegardés dans : %s\n", *outPtr)
		}

		if err := io.WriteStructured(*outPtr, structuredResult); err != nil {
			log.Fatalf("❌ Erreur lors de l'écriture : %v", err)
		}
	} else {
		fmt.Printf("\n🔄 Extraction finale des données de %s...\n", *urlPtr)
		doc, err := f.Fetch(*urlPtr)
		if err != nil {
			log.Fatalf("❌ Erreur lors de la récupération finale : %v", err)
		}

		var results []io.Result
		totalMatches := 0
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
			totalMatches += len(matches)
		}

		fmt.Printf("✅ Extraction terminée : %d sélecteurs utilisés, %d éléments extraits\n", len(selectors), totalMatches)

		if *outPtr == "-" {
			fmt.Printf("📤 Résultats affichés ci-dessous :\n\n")
		} else {
			fmt.Printf("📁 Résultats sauvegardés dans : %s\n", *outPtr)
		}

		if err := io.Write(*outPtr, io.DocumentResult{URL: *urlPtr, Results: results}); err != nil {
			log.Fatalf("❌ Erreur lors de l'écriture : %v", err)
		}
	}
}

func interactiveSession(startURL string, f *fetcher.Fetcher) ([]string, map[string]interface{}, bool, error) {
	var collectedSelectors []string
	var structuredData map[string]interface{}
	var useStructuredOutput bool
	currentURLStr := startURL

	for {
		parsedURL, err := url.Parse(currentURLStr)
		if err != nil {
			return nil, nil, false, fmt.Errorf("invalid URL '%s': %w", currentURLStr, err)
		}

		fmt.Printf("Fetching %s...\n", currentURLStr)
		doc, err := f.Fetch(currentURLStr)
		if err != nil {
			return nil, nil, false, fmt.Errorf("fetch error for %s: %w", currentURLStr, err)
		}

		res, err := tui.PromptSelectors(doc, parsedURL)
		if err != nil {
			return nil, nil, false, fmt.Errorf("TUI prompt failed: %w", err)
		}

		if len(res.Selectors) > 0 {
			if res.SelectedData != nil {
				structuredData = res.SelectedData
				useStructuredOutput = true
				fmt.Printf("✅ Données structurées sélectionnées\n")
			} else {
				collectedSelectors = append(collectedSelectors, res.Selectors...)
				fmt.Printf("✅ Added %d selectors. Total: %d\n", len(res.Selectors), len(collectedSelectors))
			}
		}

		if res.Finished {
			break
		}

		if res.NextURL != "" {
			currentURLStr = res.NextURL
		}
	}

	if useStructuredOutput {
		return nil, structuredData, true, nil
	}

	uniqueSelectors := make(map[string]struct{})
	for _, s := range collectedSelectors {
		uniqueSelectors[s] = struct{}{}
	}
	finalSelectors := make([]string, 0, len(uniqueSelectors))
	for s := range uniqueSelectors {
		finalSelectors = append(finalSelectors, s)
	}

	return finalSelectors, nil, false, nil
}
