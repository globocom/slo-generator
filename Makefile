GOLANG_LINT_COMMAND := $(shell { command -v golangci-lint; } 2>/dev/null)

test:
	go test ./...

test_generation:
	go run main.go -slo.path=slo_example.yml -rule.output /tmp/rule.example.yml
	cat /tmp/rule.example.yml

.PHONY: lint
lint:
ifndef GOLANG_LINT_COMMAND
	@echo "Command golangci-lint not found"
	@echo "Please, run the following command as sudo to install it"
	@echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.33.1"
	@exit 1
endif
	@golangci-lint run --out-format=github-actions
