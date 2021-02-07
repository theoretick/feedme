package feedparser

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

const timeoutSeconds = 30

// Latest returns the list of latest items
func Latest(maxItems int) []*gofeed.Item {
	feedUrl := os.Getenv("RSS_FEED_URL")

	items := []*gofeed.Item{}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, error := fp.ParseURLWithContext(feedUrl, ctx)
	if error != nil {
		fmt.Printf("Error parsing %v\n%v\n", feedUrl, error.Error())
		os.Exit(1)
	}

	for n, i := range feed.Items {
		if n > maxItems {
			break
		}

		items = append(items, i)
	}

	return items
}
