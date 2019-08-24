package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"genevent/gocal"

	"github.com/k0kubun/pp"
)

func main() {
	resp, err := http.Get("https://calendar.google.com/calendar/ical/cubscouts%40pack959.com/public/basic.ics")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	start, end := time.Date(2019, 8, 5, 0, 0, 0, 0, time.UTC), time.Date(2020, 6, 30, 0, 0, 0, 0, time.UTC)

	c := gocal.NewParser(resp.Body)
	c.Start, c.End = &start, &end
	c.Parse()

	pp.Print(c.Events)

	for _, e := range c.Events {
		dateStr, err := createDateString(e.StartString, e.EndString, e.Start, e.End)
		if err != nil {
			panic(err)
		}
		fmt.Printf(template, e.Summary, dateStr, e.Location, e.Description)
	}

	os.Exit(0)
}

func createDateString(startStr, endStr string, start, end *time.Time) (string, error) {
	// If only working with dates and not times, simply check if date or
	// date range
	if len(startStr) == 8 {
		startI, err := strconv.Atoi(startStr)
		if err != nil {
			return "", err
		}
		endI, err := strconv.Atoi(endStr)
		if err != nil {
			return "", nil
		}

		if endI-startI == 1 {
			return start.Format("1/2/2006"), nil
		}
		return fmt.Sprintf("%s - %s", start.Format("1/2/2006"), end.Format("1/2/2006")), nil
	}

	// Get localized datetime
	loc, _ := time.LoadLocation("America/New_York")
	locStart := start.In(loc)
	locEnd := end.In(loc)

	// If start and end are equal, just return the start datetime
	if start.Equal(*end) {
		if start.Minute() > 0 {
			return locStart.Format("1/2/2006 3:04pm"), nil
		}
		return locStart.Format("1/2/2006 3pm"), nil
	}

	// When dates are equal
	if dateEqual(locStart, locEnd) {
		dt := start.Format("1/2/2006")
		if locStart.Format("pm") == locEnd.Format("pm") {
			return fmt.Sprintf("%s %s-%s", dt, locStart.Format("3"), createTimeString(locEnd)), nil
		}
		return fmt.Sprintf("%s %s-%s", dt, createTimeString(locStart), createTimeString(locEnd)), nil
	}

	// When dates are different
	return fmt.Sprintf("%s %s - %s %s", locStart.Format("1/2/2006"), createTimeString(locStart), locEnd.Format("1/2/2006"), createTimeString(locEnd)), nil
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func createTimeString(date time.Time) string {
	if date.Minute() > 0 {
		return date.Format("3:04pm")
	}
	return date.Format("3pm")
}

var template = `
---
title: %s
date: %s
location: %s
---

%s
`
