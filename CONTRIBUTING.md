# Guide de Contribution - WebExtractor CLI

Merci de votre intérêt pour contribuer à WebExtractor CLI ! Ce guide vous aidera à participer efficacement au projet.

## 📋 Prérequis

- **Go 1.23+** installé sur votre système
- **Git** pour le contrôle de version
- **Make** pour utiliser les tâches automatisées (optionnel mais recommandé)

## 🚀 Configuration de l'environnement

1. **Fork et clone du projet**

   ```bash
   git clone https://github.com/annaelmoussa/webextractor.git
   cd webextractor
   ```

2. **Installation des outils de développement**

   ```bash
   make install-tools
   # ou manuellement :
   go install honnef.co/go/tools/cmd/staticcheck@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. **Vérification de l'installation**
   ```bash
   make check
   ```

## 🔧 Workflow de développement

### 1. Créer une branche

```bash
git checkout -b feature/nom-de-votre-feature
# ou
git checkout -b fix/description-du-bug
```

### 2. Développement

- **Respecter les règles du projet** :

  - Dépendances limitées à stdlib + `golang.org/x/net/html` + `golang.org/x/term`
  - Couverture de tests ≥ 80%
  - Code formaté avec `gofmt`
  - Zéro warning `go vet`

- **Commandes utiles** :
  ```bash
  make fmt      # Formater le code
  make test     # Lancer les tests
  make cover    # Tests avec couverture
  make lint     # Vérifications de style
  make check    # Toutes les vérifications
  ```

### 3. Tests

- **Écrire des tests** pour toute nouvelle fonctionnalité
- **Maintenir la couverture** à ≥ 80%
- **Tester manuellement** :
  ```bash
  make build
  ./webextractor -url https://example.com -sel "h1"
  ./webextractor -url https://example.com  # mode interactif
  ```

### 4. Commit et Push

```bash
git add .
git commit -m "feat: description de votre fonctionnalité"
git push origin feature/nom-de-votre-feature
```

### 5. Pull Request

1. Créez une Pull Request sur GitHub
2. Décrivez clairement vos changements
3. Liez les issues résolues (`fixes #123`)
4. Attendez la review et les checks automatiques

## 📝 Standards de code

### Messages de commit

Utilisez la convention [Conventional Commits](https://www.conventionalcommits.org/) :

- `feat:` nouvelle fonctionnalité
- `fix:` correction de bug
- `docs:` documentation
- `style:` formatage, pas de changement de logique
- `refactor:` refactoring du code
- `test:` ajout ou modification de tests
- `chore:` tâches de maintenance

### Style de code

- **Formatage** : `gofmt -s`
- **Naming** : conventions Go standards
- **Documentation** : commentaires pour les exports publics
- **Erreurs** : gestion explicite, pas de panic sauf cas extrêmes

### Architecture

```
webextractor/
├── main.go                 # Point d'entrée minimal
├── internal/
│   ├── fetcher/           # Client HTTP
│   ├── parser/            # Parsing HTML
│   ├── tui/              # Interface utilisateur
│   └── io/               # Sortie JSON
└── *_test.go             # Tests unitaires
```

## 🧪 Tests

### Structure des tests

- **Un fichier `*_test.go`** par package
- **Tests unitaires** avec mocks si nécessaire
- **Tests d'intégration** pour les flows complets
- **Benchmarks** pour les parties critiques (optionnel)

### Exemple de test

```go
func TestFetchURL(t *testing.T) {
    tests := []struct {
        name     string
        url      string
        wantErr  bool
    }{
        {"valid URL", "https://example.com", false},
        {"invalid URL", "not-a-url", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := FetchURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("FetchURL() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 🎯 Types de contributions

### 🐛 Signalement de bugs

- Utilisez les **issues GitHub**
- Décrivez le **comportement attendu vs actuel**
- Fournissez un **exemple reproductible**
- Incluez les **détails d'environnement** (OS, Go version)

### ✨ Nouvelles fonctionnalités

- **Discutez d'abord** dans une issue
- Assurez-vous que ça respecte les **règles du projet**
- Implémentez avec **tests** et **documentation**

### 📚 Documentation

- **README.md** pour les fonctionnalités utilisateur
- **Commentaires de code** pour les APIs publiques
- **CONTRIBUTING.md** pour les développeurs

### 🔧 Optimisations

- **Benchmarks** avant/après
- **Pas de micro-optimisations** sans mesures
- **Simplicité** avant performance

## 🚫 Ce qui n'est PAS accepté

- **Dépendances supplémentaires** non autorisées
- **Code généré** dynamiquement
- **Réflexion complexe** sans justification
- **Breaking changes** sans discussion préalable
- **Code non formaté** ou avec des warnings

## 📞 Aide et questions

- **Issues GitHub** pour les questions spécifiques
- **Discussions GitHub** pour les questions générales
- **Code review** collaboratif et bienveillant

## 🎉 Reconnaissance

Tous les contributeurs sont reconnus dans la section des releases et peuvent être mentionnés dans le README.

---

**Merci de votre contribution ! 🙏**
