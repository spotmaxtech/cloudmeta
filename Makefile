tagregex="v[0-9]*.[0-9]*.[0-9]*"
version=`git describe --tags --dirty --always --long --match=${tagregex} 2>/dev/null`

default:
	go build -o awsfactory/awsfactory awsfactory/awsfactory.go

fmt:
	@gofmt -s -w .

fmt-check:
	@diff=$$(gofmt -s -d .); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

docker:
	GOOS=linux GOARCH=amd64 go build -o build/awsfactory awsfactory/awsfactory.go
	cd build && docker build -t liuzoxan/cloudmeta:${version} .

docker-push:
	docker push liuzoxan/cloudmeta:${version}

clean:
	rm -f awsfactory/awsfactory build/awsfactory

.PHONY: default fmt fmt-check docker clean
