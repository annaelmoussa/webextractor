.PHONY: help build clean test lint fmt vet check cover install-tools release-local cross-build
.DEFAULT_GOAL := help

# Variables
BINARY_NAME := webextractor
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

# Couleurs pour les messages
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

## help: Affiche cette aide
help:
	@echo "WebExtractor CLI - Makefile"
	@echo ""
	@echo "Cibles disponibles:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Développement

## fmt: Formate le code Go
fmt:
	@echo "$(GREEN)Formatage du code...$(NC)"
	@gofmt -s -w .
	@echo "$(GREEN)✓ Code formaté$(NC)"

## vet: Vérifie le code avec go vet
vet:
	@echo "$(GREEN)Vérification avec go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ go vet passé$(NC)"

## lint: Lance golangci-lint
lint:
	@echo "$(GREEN)Linting avec golangci-lint...$(NC)"
	@golangci-lint run
	@echo "$(GREEN)✓ golangci-lint passé$(NC)"

## staticcheck: Lance staticcheck
staticcheck:
	@echo "$(GREEN)Analyse statique avec staticcheck...$(NC)"
	@staticcheck ./...
	@echo "$(GREEN)✓ staticcheck passé$(NC)"

## test: Lance les tests
test:
	@echo "$(GREEN)Lancement des tests...$(NC)"
	@go test -v -race ./...
	@echo "$(GREEN)✓ Tests passés$(NC)"

## cover: Lance les tests avec couverture
cover:
	@echo "$(GREEN)Tests avec couverture...$(NC)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out | grep total | awk '{print "Couverture totale: " $$3}'
	@echo "$(GREEN)✓ Tests avec couverture terminés$(NC)"

## cover-html: Génère un rapport de couverture HTML
cover-html: cover
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Rapport HTML généré: coverage.html$(NC)"

## check: Lance toutes les vérifications (fmt, vet, lint, staticcheck, test)
check: fmt vet lint staticcheck test
	@echo "$(GREEN)✓ Toutes les vérifications passées$(NC)"

##@ Build

## build: Compile le binaire
build:
	@echo "$(GREEN)Compilation en cours...$(NC)"
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "$(GREEN)✓ Binaire compilé: $(BINARY_NAME)$(NC)"

## clean: Nettoie les fichiers générés
clean:
	@echo "$(GREEN)Nettoyage...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@rm -f coverage.out coverage.html
	@rm -rf dist/
	@echo "$(GREEN)✓ Nettoyage terminé$(NC)"

## cross-build: Compile pour toutes les plateformes
cross-build: clean
	@echo "$(GREEN)Compilation multi-plateforme...$(NC)"
	@mkdir -p dist/
	
	@echo "$(YELLOW)Building for Linux AMD64...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-amd64 .
	
	@echo "$(YELLOW)Building for Linux ARM64...$(NC)"
	@GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-arm64 .
	
	@echo "$(YELLOW)Building for Darwin AMD64...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "$(YELLOW)Building for Darwin ARM64...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "$(YELLOW)Building for Windows AMD64...$(NC)"
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "$(GREEN)✓ Compilation multi-plateforme terminée$(NC)"
	@ls -la dist/

##@ Release

## release-local: Prépare une release locale avec archives
release-local: cross-build
	@echo "$(GREEN)Préparation de la release locale...$(NC)"
	@cd dist && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
	zip $(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	
	@echo "$(GREEN)Génération des checksums...$(NC)"
	@cd dist && sha256sum *.tar.gz *.zip > checksums.txt
	
	@echo "$(GREEN)✓ Release locale prête dans dist/$(NC)"
	@ls -la dist/*.tar.gz dist/*.zip dist/checksums.txt

##@ Outils

## install-tools: Installe les outils de développement
install-tools:
	@echo "$(GREEN)Installation des outils de développement...$(NC)"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Outils installés$(NC)"

## deps: Télécharge et vérifie les dépendances
deps:
	@echo "$(GREEN)Téléchargement des dépendances...$(NC)"
	@go mod download
	@go mod verify
	@echo "$(GREEN)✓ Dépendances vérifiées$(NC)"

## tidy: Nettoie go.mod et go.sum
tidy:
	@echo "$(GREEN)Nettoyage des dépendances...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ go.mod et go.sum nettoyés$(NC)"

##@ Git

## tag: Crée un tag Git (usage: make tag VERSION=v1.0.0)
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make tag VERSION=v1.0.0$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Création du tag $(VERSION)...$(NC)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "$(GREEN)✓ Tag $(VERSION) créé et pushé$(NC)"

## version: Affiche la version actuelle
version:
	@echo "Version: $(VERSION)" 