.PHONY: clean build synth deploy test

clean:
	rm -rf build
	rm -rf cdk.out

build:
	go run scripts/build/main.go

synth:
	cdk synth

deploy:
	cdk deploy

test:
	bash -c 'diff -u <(echo -n) <(go fmt $(go list ./...))'
	go vet ./...
	go test ./... -v && (echo "\nResult=OK") || (echo "\nResult=FAIL" && exit 1)
