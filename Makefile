.PHONY: build

DIST_DIR := ./.aws-sam/build

build:
	statik -src=app/pages
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/AppFunction/app ./app
	sam build
