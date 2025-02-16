.PHONY: release

release:
	@echo "Fetching GitHub token using 'gh auth token'..."
	@export GITHUB_TOKEN=$$(gh auth token | tr -d '\n'); \
	echo "GitHub token set as environment variable. $$(echo $$GITHUB_TOKEN)"; \
	echo "Please enter the tag identifier:"; \
	read -p "Tag: " TAG; \
	echo "Please enter the tag message:"; \
	read -p "Message: " MESSAGE; \
	git tag -a $$TAG -m "$$MESSAGE"; \
	git push origin $$TAG; \
	goreleaser release --clean
