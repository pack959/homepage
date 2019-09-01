HUGO=hugo
HUGOARGS?=--minify

build: genevents functions
	$(HUGO) $(HUGOARGS)

genevents:
	cd scripts/genevent && go run main.go

functions:
    mkdir -p functions
	go build -o functions/checkout ./assets/lambda/checkout

default: build

.PHONY: build genevents default
