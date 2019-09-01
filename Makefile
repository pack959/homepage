HUGO=hugo
HUGOARGS?=--minify

build: genevents functions
	$(HUGO) $(HUGOARGS)

genevents:
	cd scripts/genevent && go run main.go

functions:
	mkdir -p functions
	cd ./assets/lambda/manual-deploy && go get ./... && go build -o ../../../functions/manual-deploy .

default: build

.PHONY: build genevents functions default
