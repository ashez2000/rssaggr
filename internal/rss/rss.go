package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ashez2000/rssaggr/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchRSSFeed(url string) (RSSFeed, error) {
	rssFeed := RSSFeed{}
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Get(url)
	if err != nil {
		return rssFeed, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return rssFeed, err
	}

	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return rssFeed, err
	}

	return rssFeed, nil
}

func AggrRSSFeeds(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Println("Running AggrRSSFeeds")
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds from database", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go fetchRSSFeed(wg, db, feed)
		}
		wg.Wait()
	}
}

func fetchRSSFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	rssFeed, err := FetchRSSFeed(feed.Url)
	if err != nil {
		log.Println(err)
	}

	_, err = db.UpdateLastFetchedAt(context.Background(), database.UpdateLastFetchedAtParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time: time.Now().UTC(),
		},
	})
	if err != nil {
		log.Println("Error updating last_fetched_at", err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("error parsing published at", err)
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			PublishedAt: publishedAt,
			CreatedAt:   time.Now().UTC(),
			FeedID:      feed.ID,
		})

		if err != nil {
			// print error if its not an duplicate key error
			if !strings.Contains(err.Error(), "duplicate key") {
				log.Println("error creating post")
			}
		}

	}

	log.Printf("Feed %s fetched, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
