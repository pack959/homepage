package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"genevent/gocal"
)

var loc, _ = time.LoadLocation("America/New_York")

func main() {
	resp, err := http.Get("https://calendar.google.com/calendar/ical/cubscouts%40pack959.com/public/basic.ics")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	start, end := time.Date(2019, 8, 5, 0, 0, 0, 0, time.UTC), time.Date(2020, 7, 31, 0, 0, 0, 0, time.UTC)

	c := gocal.NewParser(resp.Body)
	c.Start, c.End = &start, &end
	c.Parse()

	for _, e := range c.Events {
		dateStr, startDate, err := createDateString(e.StartString, e.EndString, e.Start, e.End)
		if err != nil {
			panic(err)
		}

		// Clean up locations
		location := strings.Replace(e.Location, "Grace Lutheran Church, Hockessin, Delaware 19707", "Grace Lutheran Church", -1)
		location = strings.TrimSuffix(location, ", USA")

		// Replace text in the description that is incompatible with Hugo content
		description := strings.Replace(e.Description, "\\n", "<br>", -1)

		// Only include letters and numbers in url
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		sanitizedTitle := reg.ReplaceAllString(strings.TrimPrefix(e.Summary, "Tentative: "), "")
		url := fmt.Sprintf("calendar/%s_%s", e.Start.Format("20060102"), sanitizedTitle)

		data := fmt.Sprintf(template,
			e.Summary,
			startDate.Format("2006-01-02T15:04:05-0700"),
			dateStr,
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			location,
			url,
			description)
		filename := fmt.Sprintf("../../content/calendar/%s.md", strings.Replace(e.Uid, "@google.com", "", 1))
		writeToFile(filename, data)
	}

	os.Exit(0)
}

func createDateString(startStr, endStr string, start, end *time.Time) (string, *time.Time, error) {
	// If only working with dates and not times, simply check if date or
	// date range
	if len(startStr) == 8 {
		startI, err := strconv.Atoi(startStr)
		if err != nil {
			return "", nil, err
		}
		endI, err := strconv.Atoi(endStr)
		if err != nil {
			return "", nil, nil
		}

		if endI-startI == 1 {
			return start.Format("1/2/2006"), start, nil
		}
		return fmt.Sprintf("%s - %s", start.Format("1/2/2006"), end.Format("1/2/2006")), start, nil
	}

	// Get localized datetime
	loc, _ := time.LoadLocation("America/New_York")
	locStart := start.In(loc)
	locEnd := end.In(loc)

	// If start and end are equal, just return the start datetime
	if start.Equal(*end) {
		if start.Minute() > 0 {
			return locStart.Format("1/2/2006 3:04pm"), &locStart, nil
		}
		return locStart.Format("1/2/2006 3pm"), &locStart, nil
	}

	// When dates are equal
	if dateEqual(locStart, locEnd) {
		dt := locStart.Format("1/2/2006")
		if locStart.Format("pm") == locEnd.Format("pm") {
			return fmt.Sprintf("%s %s-%s", dt, createTimeString(locStart, true), createTimeString(locEnd, false)), &locStart, nil
		}
		return fmt.Sprintf("%s %s-%s", dt, createTimeString(locStart, false), createTimeString(locEnd, false)), &locStart, nil
	}

	// When dates are different
	return fmt.Sprintf("%s %s - %s %s", locStart.Format("1/2/2006"), createTimeString(locStart, false), locEnd.Format("1/2/2006"), createTimeString(locEnd, false)), &locStart, nil
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func createTimeString(date time.Time, omitAmPm bool) string {
	format := "3"
	if date.Minute() > 0 {
		format += ":04"
	}
	if !omitAmPm {
		format += "pm"
	}
	return date.Format(format)
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

var template = `---
title: "%s"
date: "%s"
dateString: "%s"
publishDate: "%s"
expiryDate: "%s"
location: "%s"
url: "%s"
---

%s
`
