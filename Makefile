BASE_DIR=$(CURDIR)

HUGO=hugo
HUGOARGS?=--minify

build: genevents
	$(HUGO) $(HUGOARGS)

genevents:
	cd scripts/genevent && go run main.go

default: build

.PHONY: build genevents default
