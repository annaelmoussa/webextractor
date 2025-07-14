package cli

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Flags contient toutes les valeurs des paramètres de ligne de commande
type Flags struct {
	URL     string
	Sel     string
	Out     string
	Timeout time.Duration
}

// Parse analyse les arguments de ligne de commande et retourne les valeurs des paramètres
func Parse() (*Flags, error) {
	flags := &Flags{
		Out:     "-",                // par défaut
		Timeout: 10 * time.Second,   // par défaut
	}
	
	args := os.Args[1:] // On ignore le nom du programme
	
	for i := 0; i < len(args); i++ {
		arg := args[i]
		
		if !strings.HasPrefix(arg, "-") {
			return nil, fmt.Errorf("unknown argument: %s", arg)
		}
		
		switch arg {
		case "-url":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-url requires a value")
			}
			flags.URL = args[i+1]
			i++ // ignore l'argument suivant (la valeur)
			
		case "-sel":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-sel requires a value")
			}
			flags.Sel = args[i+1]
			i++ // ignore l'argument suivant (la valeur)
			
		case "-out":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-out requires a value")
			}
			flags.Out = args[i+1]
			i++ // ignore l'argument suivant (la valeur)
			
		case "-timeout":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-timeout requires a value")
			}
			timeoutStr := args[i+1]
			
			// On parse la durée manuellement (cas simples)
			var duration time.Duration
			
			if strings.HasSuffix(timeoutStr, "s") {
				// On parse les secondes
				secondsStr := strings.TrimSuffix(timeoutStr, "s")
				seconds := parseInt(secondsStr)
				if seconds < 0 {
					return nil, fmt.Errorf("invalid timeout: %s", timeoutStr)
				}
				duration = time.Duration(seconds) * time.Second
			} else if strings.HasSuffix(timeoutStr, "m") {
				// On parse les minutes
				minutesStr := strings.TrimSuffix(timeoutStr, "m")
				minutes := parseInt(minutesStr)
				if minutes < 0 {
					return nil, fmt.Errorf("invalid timeout: %s", timeoutStr)
				}
				duration = time.Duration(minutes) * time.Minute
			} else {
				// On définit la durée par défaut en secondes si pas de suffixe
				seconds := parseInt(timeoutStr)
				if seconds < 0 {
					return nil, fmt.Errorf("invalid timeout: %s", timeoutStr)
				}
				duration = time.Duration(seconds) * time.Second
			}
			flags.Timeout = duration
			i++ // ignore l'argument suivant (la valeur)
			
		case "-h", "-help", "--help":
			printUsage()
			os.Exit(0)
			
		default:
			return nil, fmt.Errorf("unknown flag: %s", arg)
		}
	}
	
	// On valide les flags requis
	if flags.URL == "" {
		return nil, fmt.Errorf("required flag missing: -url")
	}
	
	return flags, nil
}

// parseInt convertit une chaîne en entier, retourne -1 en cas d'erreur.
func parseInt(s string) int {
	if s == "" {
		return -1
	}
	
	result := 0
	for _, char := range s {
		if char < '0' || char > '9' {
			return -1
		}
		digit := int(char - '0')
		result = result*10 + digit
	}
	
	return result
}

// printUsage affiche les informations d'aide.
func printUsage() {
	fmt.Printf(`Usage of %s:
  -url string
    	URL of the web page to extract from (required)
  -sel string
    	CSS-like selector (tag, .class, #id). If omitted, interactive mode starts
  -out string
    	Output JSON file path ('-' for stdout) (default "-")
  -timeout duration
    	HTTP client timeout (default 10s)
`, os.Args[0])
}