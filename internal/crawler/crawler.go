package crawler

import (
	"context"
	"net/http"
	"time"

	// NOTICE: We import our own internal packages now!
	// Replace 'your-module-name' with the name inside your go.mod file
	// e.g. "crawler-agent/internal/config"
	"crawler_agent/internal/config"
	"crawler_agent/internal/extractor"

	"github.com/mmcdole/gofeed"
)

type Job struct {
	URL         string
	LastCrawled time.Time
}

type Result struct {
	URL        string
	BodyText   string
	StatusCode int
	Error      error
}

func Fetch(ctx context.Context, job Job, cfg config.AgentConfig) Result {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", job.URL, nil)
	if err != nil {
		return Result{URL: job.URL, Error: err}
	}
	req.Header.Set("User-Agent", cfg.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return Result{URL: job.URL, Error: err}
	}
	defer resp.Body.Close()

	// Call our new package
	cleanText, err := extractor.ExtractText(resp.Body)

	if err != nil {
		return Result{
			URL:        job.URL,
			StatusCode: resp.StatusCode,
			Error:      err,
		}
	}

	return Result{
		URL:        job.URL,
		BodyText:   cleanText,
		StatusCode: resp.StatusCode,
		Error:      nil,
	}
}

func GetLatestLinks(feedURL string) ([]string, error) {
	// Implementation to fetch latest links from a feed (e.g., RSS)
	// This is a placeholder implementation
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)

	if err != nil {
		return nil, err
	}

	var links []string
	limit := 5

	if len(feed.Items) < limit {
		limit = len(feed.Items)
	}

	for i := 0; i < limit; i++ {
		links = append(links, feed.Items[i].Link)
	}

	return links, nil
}
