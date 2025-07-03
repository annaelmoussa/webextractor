# WebExtractor CLI

[![CI/CD Pipeline](https://github.com/annaelmoussa/webextractor/actions/workflows/ci.yml/badge.svg)](https://github.com/annaelmoussa/webextractor/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/annaelmoussa/webextractor)](https://goreportcard.com/report/github.com/annaelmoussa/webextractor)
[![codecov](https://codecov.io/gh/annaelmoussa/webextractor/branch/main/graph/badge.svg)](https://codecov.io/gh/annaelmoussa/webextractor)
[![Release](https://img.shields.io/github/release/annaelmoussa/webextractor.svg)](https://github.com/annaelmoussa/webextractor/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Un extracteur HTML minimaliste écrit en Go avec une architecture modulaire et des dépendances minimales.

## 🚀 Fonctionnalités

- **Extraction HTTP/HTTPS** : Récupère n'importe quelle URL avec un User-Agent personnalisé (`WebExtractor/0.1`)
- **Sélecteurs légers** : Syntaxe CSS simplifiée sans dépendances externes
  - `tag` — nom d'élément (`p`, `div`, `span`, etc.)
  - `.class` — nom de classe (`.note`, `.content`)
  - `#id` — attribut `id` exact (`#main`, `#header`)
- **Mode interactif avancé** : Interface TUI intuitive avec affichage structuré
  - **Affichage avec emojis** : Interface claire et colorée (📄 Page, 🌐 Titre, 🔠 H1, 📝 Paragraphes, 🔗 Liens, etc.)
  - **Sélection granulaire** : Choix d'éléments individuels par indices numériques
  - **Sélections multiples** : Support des indices (`0,2,4`) et plages (`1-3`, `0,2-5,8`)
  - **Sélections personnalisées** : Combinaison libre d'éléments de différentes catégories
  - **Navigation web** : Possibilité de suivre les liens détectés pour explorer d'autres pages
  - **Aperçu en temps réel** : Prévisualisation des sélections avant extraction finale
- **Sortie JSON flexible** : Format standardisé ou structuré selon le mode utilisé
- **Qualité** : >80% de couverture de tests, zéro warning `go vet`, code `go fmt` compliant

## 📋 Formats de sortie

### Mode sélecteurs classiques

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

### Mode interactif structuré

```json
{
  "url": "https://example.com",
  "title": "Example Domain",
  "h1": ["Example Domain"],
  "paragraphs": [
    "This domain is for use in illustrative examples...",
    "More information..."
  ],
  "links": ["https://www.iana.org/domains/example"]
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
```

**Interface d'exemple** :

```
📄 Page: https://example.com
🌐 Title: Example Domain

Éléments disponibles:
✅ [ 0] 🌐 TITLE Example Domain
✅ [ 1] 🔠 H1 Example Domain
✅ [ 2] 📝 P This domain is for use in illustrative...
✅ [ 3] 📝 P More information...
✅ [ 4] 🔗 LINK More information... (https://...)

Commandes disponibles:
• Sélection: tapez les numéros (ex: 0,2,4 ou 1-3 ou all)
• Navigation: L<numéro> pour suivre un lien
• Aperçu: preview pour voir la sélection actuelle
• Sortie: done pour générer le JSON final
```

**Exemples de sélections** :

- `all` — Sélectionner tous les éléments
- `0` — Sélectionner uniquement l'élément 0 (titre)
- `1,3,4` — Sélectionner les éléments 1, 3 et 4
- `0-2` — Sélectionner les éléments 0, 1 et 2
- `0,2-4,7` — Combinaison : éléments 0, 2 à 4, et 7
- `L4` — Naviguer vers le lien de l'élément 4

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

## 🎯 Fonctionnalités avancées

### Interface TUI granulaire

- **Sélection fine** : Choisissez exactement les éléments que vous voulez
- **Aperçu instantané** : Voyez votre sélection avant l'extraction
- **Navigation intuitive** : Explorez les liens directement depuis l'interface
- **Affichage structuré** : Emojis et catégorisation pour une meilleure lisibilité

### Flexibilité d'extraction

- **Mode sélecteurs** : Pour les utilisations scriptées/automatisées
- **Mode interactif** : Pour l'exploration et la sélection précise
- **Sortie adaptative** : JSON classique ou structuré selon le contexte

## 🧪 Développement

### Pipeline CI/CD

Le projet utilise GitHub Actions pour automatiser :

- **Tests et qualité** : Exécution automatique des tests avec couverture ≥ 80%
- **Linting** : Vérification du code avec `go vet`, `gofmt`, `golangci-lint` et `staticcheck`
- **Build multi-plateforme** : Compilation pour Linux, macOS et Windows (AMD64 + ARM64)
- **Scans de sécurité** : Analyse avec `gosec` et review des dépendances
- **Releases automatiques** : Génération de binaires et archives lors des tags

### Commandes de développement

```bash
# Installation des outils de dev
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Tests avec couverture
go test ./... -cover -race -covermode=atomic

# Vérification complète (comme en CI)
gofmt -s -l . && \
go vet ./... && \
staticcheck ./... && \
golangci-lint run && \
go test ./... -cover

# Build multi-plateforme
GOOS=linux GOARCH=amd64 go build -o webextractor-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o webextractor-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o webextractor-windows-amd64.exe
```

### Tests manuels

```bash
# Test basique
./webextractor -url https://httpbin.org/html -sel "h1"

# Test mode interactif
./webextractor -url https://example.com

# Test avec fichier de sortie
./webextractor -url https://httpbin.org/html -sel "p" -out /tmp/test.json

# Test sélections granulaires
echo "0,2-4" | ./webextractor -url https://example.com

# Test navigation
echo "L0" | ./webextractor -url https://example.com
```

### Processus de release

1. **Développement** : Travail sur une branche feature
2. **Pull Request** : Tests automatiques et review
3. **Merge** : Intégration en `main` après validation
4. **Tag** : Création d'un tag `v1.2.3` pour déclencher la release
5. **Release automatique** : Génération des binaires et publication

```bash
# Créer une release
git tag v1.0.0
git push origin v1.0.0
# → Déclenche automatiquement la création de la release avec binaires
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
