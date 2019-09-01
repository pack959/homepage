HUGO=hugo
HUGOARGS?=--minify

build: genevents functions
	$(HUGO) $(HUGOARGS)

genevents:
	cd scripts/genevent && go run main.go

functions:
	mkdir -p functions
	cd ./assets/lambda/checkout && go build -o ../../../functions/checkout .

default: build

.PHONY: build genevents default functions
