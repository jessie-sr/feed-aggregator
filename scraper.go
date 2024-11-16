package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jessie-sr/rss-aggregator/internal/db"
)

func startScraping(DB *db.Queries,
	concurrency int, // # of different goroutines to run for the scraping process.
	timeBetweenRequest time.Duration, // Time interval between scraping requests
) {
	log.Printf("Scraping on %v Goroutines every %v durations", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C { // Block until the ticker sends a value on its ticker.C (happens at every timeBetweenRequest interval)
		feeds, err := DB.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue // Should always be running as our server operates
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(DB, wg, feed) // Concurrent scraping
		}
		wg.Wait() // Wait for {concurrency} number of goroutines to finish their scraping
	}
}

func scrapeFeed(DB *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	_, err := DB.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed from URL: ", err)
		return
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds) // Use precise timestamps to avoid logs being out of order due to concurrency or buffer issues

	for _, item := range rssFeed.Channel.Item {
		// log.Println("Found post", item.Title, "on feed", feed.Name) // This log is for testing purpose

		link := item.Link
		if link == "" {
			link = item.Guid // Use GUID if Link is empty
		}
		// log.Printf("Parsed link: %s", link)

		// Initialize description
		description := sql.NullString{}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// Parse time
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing date %v with error %v: ", item.PubDate, err)
			continue
		}

		_, err = DB.CreatePost(context.Background(), db.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         link,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Println("Error creating post: ", err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
