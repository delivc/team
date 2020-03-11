.PHONY: all build deps image lint migrate test race msan vet
CHECK_FILES?=$$(go list ./... | grep -v /vendor/)

lint: ## Lint the code.
	@golint -set_exit_status $(CHECK_FILES)