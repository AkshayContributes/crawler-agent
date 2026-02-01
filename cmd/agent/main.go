package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"crawler_agent/internal/brain"
	"crawler_agent/internal/config"
	"crawler_agent/internal/crawler"
	"crawler_agent/internal/notifier"
	"crawler_agent/internal/state"
)

func main() {
	cfg := config.AgentConfig{
		ConcurrencyLimit:  3,
		UserAgent:         "Mozilla/5.0 (compatible; CrawlerAgent/1.0; +http://example.com/bot)",
		Timeout:           60 * time.Second,
		RequestsPerMinute: 4,
	}

	tgBot := notifier.NewTelegramClient(os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID"))

	var brainClient brain.AI
	var err error

	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey != "" {
		modelName := os.Getenv("GEMINI_MODEL_NAME")
		brainClient, err = brain.NewGeminiStrategy(apiKey, modelName)
		if err != nil {
			fmt.Printf("Error initializing Gemini strategy: %v\n", err)
		}
	} else {
		brainClient = brain.NewOllamaStrategy("llama3.1:8b")
	}

	if err != nil {
		fmt.Printf("Error initializing brain strategy: %v\n", err)
		return
	}

	defer brainClient.Close()

	history, err := state.NewStore("history.json")

	if err != nil {
		fmt.Printf("Error initializing state store: %v\n", err)
		return
	}

	feeds := []string{
		"https://netflixtechblog.com/feed",
		"https://eng.uber.com/feed/",
		"https://go.dev/blog/feed.atom",
		"https://aws.amazon.com/blogs/architecture/feed/",
	}

	fmt.Println("🤖 Agent started. Running continuously...")

	runCycle(feeds, history, tgBot, cfg, brainClient)

	ticker := time.NewTicker(6 * time.Hour)

	for range ticker.C {
		fmt.Println("\n⏰ Waking up for scheduled scan...")
		runCycle(feeds, history, tgBot, cfg, brainClient)
	}
}
func runCycle(feeds []string, history *state.Store, tgBot *notifier.TelegramClient, cfg config.AgentConfig, brainClient brain.AI) {

	safeDelay := time.Minute / time.Duration(cfg.RequestsPerMinute)

	var newLinks []string
	for _, feed := range feeds {
		links, err := crawler.GetLatestLinks(feed)
		if err != nil {
			fmt.Printf("Error fetching links from feed %s: %v\n", feed, err)
			continue
		}

		for _, link := range links {
			if !history.HasSeen(link) {
				newLinks = append(newLinks, link)
			}
		}
	}

	if len(newLinks) == 0 {
		fmt.Println("💤 No new articles found.")
		return
	}

	fmt.Printf("🚀 Found %d NEW articles. Starting analysis...\n", len(newLinks))

	resultsChan := make(chan crawler.Result, len(newLinks))

	sem := make(chan struct{}, cfg.ConcurrencyLimit)

	var wg sync.WaitGroup

	for _, url := range newLinks {
		wg.Add(1)
		sem <- struct{}{}

		go func(targetUrl string) {
			defer wg.Done()
			defer func() { <-sem }()

			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()

			fmt.Printf("🔍 Processing: %s\n", targetUrl)

			job := crawler.Job{URL: targetUrl}

			res := crawler.Fetch(ctx, job, cfg)

			if res.Error != nil {
				fmt.Printf("❌ Error processing %s: %v\n", targetUrl, res.Error)
			} else {
				summary, err := brainClient.Summarize(res.BodyText)
				if err != nil {
					fmt.Printf("❌ Error summarizing %s: %v\n", targetUrl, err)
				} else {
					res.BodyText = summary
					fmt.Printf("✅ Completed: %s\n", targetUrl)
				}
			}
			resultsChan <- res
		}(url)

		time.Sleep(safeDelay)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		if res.Error != nil {
			fmt.Printf("❌ Failed: %s\n", res.URL)
			continue
		}

		msg := fmt.Sprintf("📰 *New Article Analyzed*\n\n*URL:* %s\n*Summary:*\n%s", res.URL, res.BodyText)
		if err := tgBot.SendMessage(msg); err != nil {
			fmt.Printf("Error sending Telegram message for %s: %v\n", res.URL, err)
		} else {
			fmt.Printf("📨 Notification sent for %s\n", res.URL)
		}

		if err := history.MarkSeen(res.URL); err != nil {
			fmt.Printf("Error updating history for %s: %v\n", res.URL, err)
		}
	}
}
