HUGO=hugo
HUGOARGS?=--minify --gc

build: buildgo genevents
	$(HUGO) $(HUGOARGS)

genevents: buildgo
	./bin/genevent https://calendar.google.com/calendar/ical/cubscouts%40pack959.com/public/basic.ics 8/5/2019 7/31/2020 ./content/calendar/

buildgo:
	mkdir -p bin
	cd scripts/genevent && go build -o ../../bin/genevent .

default: build

.PHONY: build buildgo genevents default
