HUGO=hugo
HUGOARGS?=--minify --gc

build: functions genevents
	$(HUGO) $(HUGOARGS)

genevents: buildgo
	./bin/genevent https://calendar.google.com/calendar/ical/cubscouts%40pack959.com/public/basic.ics 8/5/2019 7/31/2020 ./content/calendar/

buildgo:
	mkdir -p bin
	cd scripts/genevent && go build -o ../../bin/genevent .

functions:
	mkdir -p functions
	cd ./assets/lambda/manual-deploy && go build -o ../../../functions/manual-deploy .

default: build

.PHONY: build buildgo genevents functions default
