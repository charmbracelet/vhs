themes.json:
	# this is the url used by https://windowsterminalthemes.dev/
	# See https://github.com/atomcorp/themes/blob/master/app/src/App.tsx#L18
	@./scripts/download_theme.sh \
		https://2zrysvpla9.execute-api.eu-west-2.amazonaws.com/prod/themes \
		themes

THEMES.md:
	@go run . themes --markdown 2> THEMES.md

all: themes.json THEMES.md
	@echo "Running all"

refresh:
	@rm -rf themes.json themes.json THEMES.md
	@$(MAKE) all

