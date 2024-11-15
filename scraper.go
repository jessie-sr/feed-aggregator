package main

import (
	"log"
	"time"

	"github.com/jessie-sr/rss-aggregator/internal/db"
)

func startScraping(db *db.Queries,
	concurrency int, // # of different goroutines to run for the scraping process.
	timeBetweenRequest time.Duration, // Time interval between scraping requests
) {
	log.Printf("Scraping on %v Goroutines every %v durations", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C { // Block until the ticker sends a value on its ticker.C (happens at every timeBetweenRequest interval)

	}
}
