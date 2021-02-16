PRECOMMIT_COMMAND := $(shell { command -v pre-commit; } 2>/dev/null)
GOLANG_LINT_COMMAND := $(shell { command -v golangci-lint; } 2>/dev/null)

test:
	go test ./...

test_generation:
	go run main.go -slo.path=slo_example.yml -rule.output /tmp/rule.example.yml
	cat /tmp/rule.example.yml

setup: pre-commit
	go get .

.PHONY: pre-commit
pre-commit:
ifndef PRECOMMIT_COMMAND
	@echo "\nCommand 'pre-commit' not found!\n"
	@echo "Please, run the following command to install it:"
	@echo "\nMacOSX:"
	@echo "brew install pre-commit"
	@echo "\nGNU/Linux:"
	@echo "aptitude install pre-commit"
	@echo "\nMore info, take a look at: https://pre-commit.com/index.html#install\n\n"
	@exit 1
endif
	@pre-commit install --install-hooks
	@pre-commit install --hook-type pre-push

.PHONY: lint
lint:
ifndef GOLANG_LINT_COMMAND
	@echo "Command golangci-lint not found"
	@echo "Please, run the following command as sudo to install it"
	@echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.33.1"
	@exit 1
endif
	@golangci-lint run --out-format=github-actions