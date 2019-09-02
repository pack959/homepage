HUGO=hugo
HUGOARGS?=--minify --gc
CALSTART=8/5/2019
CALEND=12/31/2020

build: functions genevents
	$(HUGO) $(HUGOARGS)

genevents: buildgo
	./bin/genevent https://calendar.google.com/calendar/ical/cubscouts%40pack959.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/
	./bin/genevent https://calendar.google.com/calendar/ical/pack959.com_lmh78ccr9gdv49573auvhbh21s%40group.calendar.google.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/lion/
	./bin/genevent https://calendar.google.com/calendar/ical/pack959.com_n9ivp4528uctq7n8somko9e1d0%40group.calendar.google.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/tiger/
	./bin/genevent https://calendar.google.com/calendar/ical/pack959.com_v53027ab8vgo1jl78ko78aocfg%40group.calendar.google.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/wolf/
	./bin/genevent https://calendar.google.com/calendar/ical/pack959.com_b5633aklqu4e31cldfou9509gs%40group.calendar.google.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/bear/
	./bin/genevent https://calendar.google.com/calendar/ical/pack959.com_3pjo65dmf1oh075suaspf04kt4%40group.calendar.google.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/webelos/
	./bin/genevent https://calendar.google.com/calendar/ical/pack959.com_0i2jo8ahu9stn0o9ngboh6ivhg%40group.calendar.google.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/aol/	
	./bin/genevent https://calendar.google.com/calendar/ical/pack959.com_qarlih8q0akjb5otvt7aclh2hg%40group.calendar.google.com/public/basic.ics $(CALSTART) $(CALEND) ./content/calendar/leader/	

buildgo:
	mkdir -p bin
	cd scripts/genevent && go build -o ../../bin/genevent .

functions:
	mkdir -p functions
	cd ./assets/lambda/manual-deploy && go build -o ../../../functions/manual-deploy .

default: build

.PHONY: build buildgo genevents functions default
