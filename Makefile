themes.json:
	@./scripts/download_theme.sh backupthemes themes

themes_custom.json:
	@./scripts/download_theme.sh custom-colour-schemes themes_custom

THEMES.md:
	@go run . themes --markdown > THEMES.md

all: themes.json themes_custom.json THEMES.md
	@echo "Running all"

refresh:
	@rm -rf themes.json themes.json THEMES.md
	@$(MAKE) all

