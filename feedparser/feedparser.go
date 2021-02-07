package feedparser

import (
	"os"

	"github.com/mmcdole/gofeed"
)

// Latest returns the list of latest items
func Latest(maxItems int) []*gofeed.Item {
	feedUrl := os.Getenv("RSS_FEED_URL")

	items := []*gofeed.Item{}
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedUrl)

	for n, i := range feed.Items {
		if n > maxItems {
			break
		}

		items = append(items, i)
	}

	return items
}
