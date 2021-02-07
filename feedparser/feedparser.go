package feedparser

import (
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
)

// Latest returns the list of latest items
func Latest(maxItems int) []*gofeed.Item {
	feedUrl := os.Getenv("RSS_FEED_URL")

	items := []*gofeed.Item{}
	fp := gofeed.NewParser()
	feed, error := fp.ParseURL(feedUrl)
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
