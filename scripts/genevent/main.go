package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/apognu/gocal"
)

func main() {
    resp, err := http.Get("https://calendar.google.com/calendar/ical/cubscouts%40pack959.com/public/basic.ics")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
 
	start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)
  
	c := gocal.NewParser(resp.Body)
	c.Start, c.End = &start, &end
	c.Parse()
  
	for _, e := range c.Events {
	  fmt.Printf("%s\n", e.Summary)
	}  

    os.Exit(42)
}