package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

type ChannelInfo struct {
	ID       string
	Name     string
	LastSeen time.Time
}

type SkipList struct {
	SkipChannels []string `json:"skip_channels"`
}

var skipChannels = map[string]bool{}

func initSkipList() {
	// Try to read from JSON file first
	data, err := os.ReadFile("skip-channels.json")
	if err != nil {
		log.Printf("⚠️ Warning: Could not read skip-channels.json, using default skip list: %v", err)
		// Fallback to default skip list
		defaultSkipList := []string{
			"general",
			"company-announcements",
			"team-leads",
			"support-team",
			"dev-team",
			"infrastructure",
			"security",
			"hr-announcements",
			"product-updates",
			"sales-team",
			"marketing-team",
			"customer-support",
			"emergency-alerts",
			"system-notifications",
			"admin-only",
			"executive-team",
			"board-meetings",
			"legal-team",
			"finance-team",
			"compliance",
		}
		for _, name := range defaultSkipList {
			skipChannels[name] = true
		}
		return
	}

	var skipList SkipList
	if err := json.Unmarshal(data, &skipList); err != nil {
		log.Printf("⚠️ Warning: Could not parse skip-channels.json, using default skip list: %v", err)
		// Fallback to default skip list
		defaultSkipList := []string{
			"general",
			"company-announcements",
			"team-leads",
			"support-team",
			"dev-team",
			"infrastructure",
			"security",
			"hr-announcements",
			"product-updates",
			"sales-team",
			"marketing-team",
			"customer-support",
			"emergency-alerts",
			"system-notifications",
			"admin-only",
			"executive-team",
			"board-meetings",
			"legal-team",
			"finance-team",
			"compliance",
		}
		for _, name := range defaultSkipList {
			skipChannels[name] = true
		}
		return
	}

	// Load channels from JSON file
	for _, name := range skipList.SkipChannels {
		skipChannels[strings.TrimSpace(name)] = true
	}
	
	if len(skipChannels) > 0 {
		log.Printf("✅ Loaded %d channels to skip from skip-channels.json", len(skipChannels))
	}
}

func main() {
	_ = godotenv.Load()
	initSkipList()

	days := flag.Int("days", 0, "Filter channels with last activity older than this many days")
	keyword := flag.String("keyword", "", "Filter channels whose names contain this keyword")
	channelTypes := flag.String("types", "public,private", "Channel types to include: public,private or both (comma-separated)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *days == 0 && *keyword == "" {
		log.Fatal("❌ You must provide either --days or --keyword")
	}

	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		log.Fatal("❌ SLACK_API_TOKEN not set in environment or .env file")
	}

	api := slack.New(token)
	cutoff := time.Now().AddDate(0, 0, -*days)

	var types []string
	for _, t := range strings.Split(*channelTypes, ",") {
		t = strings.TrimSpace(strings.ToLower(t))
		if t == "public" {
			types = append(types, "public_channel")
		} else if t == "private" {
			types = append(types, "private_channel")
		}
	}
	if len(types) == 0 {
		types = []string{"public_channel", "private_channel"}
	}

	if *days > 0 && *verbose {
		log.Printf("⏳ Checking for stale channels older than %s...\n", cutoff.Format(time.RFC1123))
	}
	if *keyword != "" && *verbose {
		log.Printf("🔍 Filtering channels with keyword: %q\n", *keyword)
	}
	if *verbose {
		log.Printf("📡 Including channel types: %v\n", types)
	}

	filteredChannels, err := getFilteredChannels(api, cutoff, *keyword, *days > 0, types)
	if err != nil {
		log.Fatalf("❌ Error getting channels: %v\n", err)
	}

	if len(filteredChannels) == 0 {
		log.Println("🎉 No matching channels found.")
		return
	}

	log.Printf("📦 Found %d channel(s) matching criteria.\n", len(filteredChannels))
	for i, ch := range filteredChannels {
		fmt.Printf("[%d/%d] #%-30s ID: %s", i+1, len(filteredChannels), ch.Name, ch.ID)
		if !ch.LastSeen.IsZero() {
			fmt.Printf(" — Last message: %s", ch.LastSeen.Format(time.RFC1123))
		}
		fmt.Println()
	}

	fmt.Print("\n❓ Do you want to leave these channels now? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		leaveChannels(api, filteredChannels)
	} else {
		log.Println("👍 Skipping leave. No channels were left.")
	}
}

func leaveChannels(api *slack.Client, channels []ChannelInfo) {
	log.Println("🚪 Leaving channels...")
	for i, ch := range channels {
		log.Printf("➡️  [%d/%d] Leaving #%s (ID: %s)...", i+1, len(channels), ch.Name, ch.ID)
		_, err := api.LeaveConversation(ch.ID)
		if err != nil {
			log.Printf("❌ Failed to leave #%s: %v", ch.Name, err)
		} else {
			log.Printf("✅ Left #%s", ch.Name)
		}
		time.Sleep(1 * time.Second)
	}
}

func getFilteredChannels(api *slack.Client, cutoff time.Time, keyword string, useDate bool, types []string) ([]ChannelInfo, error) {
	var results []ChannelInfo
	var wg sync.WaitGroup
	chMutex := sync.Mutex{}
	semaphore := make(chan struct{}, 5)
	cursor := ""

	for {
		var channels []slack.Channel
		var nextCursor string
		var err error

		for {
			params := &slack.GetConversationsParameters{
				Limit:           30,
				ExcludeArchived: true,
				Cursor:          cursor,
				Types:           types,
			}

			channels, nextCursor, err = api.GetConversations(params)
			if err != nil {
				if strings.Contains(err.Error(), "rate_limited") {
					log.Println("⚠️ Hit rate limit. Waiting 30s...")
					time.Sleep(30 * time.Second)
					continue
				}
				if rateErr, ok := err.(*slack.RateLimitedError); ok {
					wait := time.Duration(rateErr.RetryAfter) * time.Second
					if wait <= 0 {
						wait = 30 * time.Second
					}
					log.Printf("⏳ Rate limit hit while fetching channels. Waiting %v before retrying...", wait)
					time.Sleep(wait)
					continue
				}
				return nil, fmt.Errorf("error listing conversations: %w", err)
			}
			break
		}

		for _, ch := range channels {
			if !ch.IsMember || skipChannels[ch.Name] {
				continue
			}
			if keyword != "" && !strings.Contains(ch.Name, keyword) {
				continue
			}

			wg.Add(1)
			go func(ch slack.Channel) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				var lastTime time.Time
				if useDate {
					var history *slack.GetConversationHistoryResponse
					for attempt := 0; attempt < 2; attempt++ {
						history, err = api.GetConversationHistory(&slack.GetConversationHistoryParameters{
							ChannelID: ch.ID,
							Limit:     1,
						})
						if err != nil {
							if strings.Contains(err.Error(), "rate_limited") {
								time.Sleep(30 * time.Second)
								continue
							}
							if rateErr, ok := err.(*slack.RateLimitedError); ok {
								wait := time.Duration(rateErr.RetryAfter) * time.Second
								if wait <= 0 {
									wait = 30 * time.Second
								}
								time.Sleep(wait)
								continue
							}
							return
						}
						break
					}
					if history == nil || len(history.Messages) == 0 {
						return
					}

					tsFloat, err := strconv.ParseFloat(history.Messages[0].Timestamp, 64)
					if err != nil {
						return
					}
					lastTime = time.Unix(int64(tsFloat), 0)

					if lastTime.After(cutoff) {
						return
					}
				}

				chMutex.Lock()
				results = append(results, ChannelInfo{
					ID:       ch.ID,
					Name:     ch.Name,
					LastSeen: lastTime,
				})
				chMutex.Unlock()

				time.Sleep(1 * time.Second)
			}(ch)
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	wg.Wait()
	return results, nil
} 