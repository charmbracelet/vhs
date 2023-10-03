themes.json:
	@./scripts/download_theme.sh \
		https://2zrysvpla9.execute-api.eu-west-2.amazonaws.com/prod/themes \
		themes

themes_custom.json:
	@./scripts/download_theme.sh \
		https://raw.githubusercontent.com/atomcorp/themes/master/app/src/custom-colour-schemes.json \
		themes_custom

THEMES.md:
	@go run . themes --markdown 2> THEMES.md

all: themes.json themes_custom.json THEMES.md
	@echo "Running all"

refresh:
	@rm -rf themes.json themes.json THEMES.md
	@$(MAKE) all

