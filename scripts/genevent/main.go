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

	"github.com/apognu/gocal"
)

var loc, _ = time.LoadLocation("America/New_York")

func main() {
	args := os.Args[1:]

	resp, err := http.Get(args[0])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	start, err := time.Parse("1/2/2006", args[1])
	if err != nil {
		log.Fatalf("error parsing start date: %s", err)
	}

	end, err := time.Parse("1/2/2006", args[2])
	if err != nil {
		log.Fatalf("error parsing start date: %s", err)
	}

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
		slug := fmt.Sprintf("%s_%s", e.Start.Format("20060102"), sanitizedTitle)

		data := fmt.Sprintf(template,
			e.Summary,
			startDate.Format("2006-01-02T15:04:05-0700"),
			dateStr,
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			location,
			slug,
			description)
		filename := fmt.Sprintf("%s/%s.md", args[3], strings.Replace(e.Uid, "@google.com", "", 1))
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

		locStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
		if endI-startI == 1 {
			// 1/1/2019
			return start.Format("1/2/2006"), &locStart, nil
		}
		// 1/1/2019 - 1/3/2019
		return fmt.Sprintf("%s - %s", start.Format("1/2/2006"), end.Format("1/2/2006")), &locStart, nil
	}

	locStart := start.In(loc)
	locEnd := end.In(loc)

	// If start and end are equal, just return the start datetime
	if start.Equal(*end) {
		// 1/1/2019 7:30pm
		return fmt.Sprintf("%s %s", locStart.Format("1/2/2006"), createTimeString(locStart, false)), &locStart, nil
	}

	// When dates are equal
	if dateEqual(locStart, locEnd) {
		dt := locStart.Format("1/2/2006")
		if locStart.Format("pm") == locEnd.Format("pm") {
			// 1/1/2019 1-2pm
			return fmt.Sprintf("%s %s-%s", dt, createTimeString(locStart, true), createTimeString(locEnd, false)), &locStart, nil
		}
		// 1/1/2019 1am-2pm
		return fmt.Sprintf("%s %s-%s", dt, createTimeString(locStart, false), createTimeString(locEnd, false)), &locStart, nil
	}

	// When dates are different
	// 7/1/2019 8am - 7/4/2019 4:30pm
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
slug: "%s"
---

%s
`
