# WebExtractor CLI

Un extracteur HTML minimaliste écrit en Go avec une architecture modulaire et des dépendances minimales.

## 🚀 Fonctionnalités

- **Extraction HTTP/HTTPS** : Récupère n'importe quelle URL avec un User-Agent personnalisé (`WebExtractor/0.1`)
- **Sélecteurs légers** : Syntaxe CSS simplifiée sans dépendances externes
  - `tag` — nom d'élément (`p`, `div`, `span`, etc.)
  - `.class` — nom de classe (`.note`, `.content`)
  - `#id` — attribut `id` exact (`#main`, `#header`)
- **Mode interactif** : Interface TUI intuitive quand `-sel` est omis
  - Navigation par pages (15 éléments max par page)
  - Sélection par indices numériques (`0`, `1,3,5`, `0-3`)
  - Catégorisation automatique des éléments (Titres, Textes, Liens, etc.)
  - Navigation entre pages web
- **Sortie JSON** : Format standardisé avec indentation 2 espaces
- **Qualité** : >80% de couverture de tests, zéro warning `go vet`, code `go fmt` compliant

## 📋 Format de sortie

```json
{
  "url": "https://example.com",
  "results": [
    {
      "selector": "h1",
      "matches": ["Titre principal"]
    },
    {
      "selector": ".note",
      "matches": ["Note importante", "Autre note"]
    }
  ]
}
```

## 🛠 Installation

```bash
# Cloner et compiler
cd webextractor
go build
```

## 💡 Utilisation

### Mode direct (avec sélecteurs)

```bash
# Extraction simple
./webextractor -url https://example.com -sel "h1"

# Plusieurs sélecteurs
./webextractor -url https://example.com -sel "h1,p,.note" -out data.json

# Avec timeout personnalisé
./webextractor -url https://example.com -sel "#main" -timeout 30s
```

### Mode interactif (sans sélecteurs)

```bash
# Lancer le mode interactif
./webextractor -url https://example.com

# L'interface vous guidera pour :
# - Voir les éléments par catégories
# - Sélectionner par numéros (0, 1,3,5, 0-3)
# - Naviguer vers d'autres pages (L0, L1, etc.)
# - Prévisualiser l'extraction avant de terminer
```

### Paramètres disponibles

| Paramètre  | Description                         | Défaut          |
| ---------- | ----------------------------------- | --------------- |
| `-url`     | URL cible **(requis)**              | -               |
| `-sel`     | Sélecteurs CSS séparés par virgules | Mode interactif |
| `-out`     | Chemin de sortie (`-` pour stdout)  | `-`             |
| `-timeout` | Timeout HTTP                        | `10s`           |

## 🏗 Architecture

```
webextractor/
├── main.go                 # Point d'entrée et orchestration
├── internal/
│   ├── fetcher/           # Client HTTP avec timeout
│   ├── parser/            # Analyseur HTML et sélecteurs CSS
│   ├── tui/              # Interface utilisateur interactive
│   └── io/               # Sortie JSON formatée
└── *_test.go             # Tests unitaires (>80% couverture)
```

### Principes de conception

- **Dépendances minimales** : Seulement stdlib + `golang.org/x/net/html`
- **Modules séparés** : Chaque package a une responsabilité unique
- **API simple** : Pas de réflexion complexe ou de code généré
- **Testabilité** : Chaque composant est testé unitairement

## 🧪 Développement

### Tests et qualité

```bash
# Tests avec couverture
go test ./... -cover

# Vérification statique
go vet ./...

# Formatage du code
go fmt ./...

# Pipeline complet
go vet ./... && go fmt ./... && go test ./... -cover
```

### Tests manuels

```bash
# Test basique
./webextractor -url https://httpbin.org/html -sel "h1"

# Test mode interactif
./webextractor -url https://example.com

# Test avec fichier de sortie
./webextractor -url https://httpbin.org/html -sel "p" -out /tmp/test.json
```

## 📦 Dépendances

- **Go 1.22+** (standard library)
- `golang.org/x/net/html` — Parseur HTML
- `golang.org/x/term` — Support terminal (optionnel, pour futures améliorations TUI)

## 🤝 Contribution

1. Fork le projet
2. Créer une branche (`git checkout -b feature/ma-fonctionnalite`)
3. Commit les changements (`git commit -am 'Ajouter ma fonctionnalité'`)
4. Push vers la branche (`git push origin feature/ma-fonctionnalite`)
5. Ouvrir une Pull Request

## 📄 Licence

Ce projet respecte les règles d'architecture définies pour un CLI Go minimal et efficace.
