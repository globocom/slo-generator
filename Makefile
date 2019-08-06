test:
	go test ./...

test_generation:
	go run main.go -slo.path=slo_example.yml -rule.output /tmp/rule.example.yml
	cat /tmp/rule.example.yml
