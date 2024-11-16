package main

import (
	"context"
	"log"
	"sync"
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
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue // Should always be running as our server operates
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed) // Concurrent scraping
		}
		wg.Wait() // Wait for {concurrency} number of goroutines to finish their scraping
	}
}

func scrapeFeed(db *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed from URL: ", err)
		return
	}

	// This log is for testing purpose
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) // Use precise timestamps to avoid logs being out of order due to concurrency or buffer issues

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
