.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/create lambdas/create/create.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/delete lambdas/delete/delete.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get lambdas/get/get.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/list lambdas/list/list.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update lambdas/update/update.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
