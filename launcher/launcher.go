package launcher

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
)

// OpenUrl opens an Item.Link using the host machine's default
// browser. There's probably a lot of better ways to do this.
func OpenUrl(rel *gofeed.Item) {
	url := rel.Link
	var downloadLink = regexp.MustCompile(`action=download`)
	if downloadLink.MatchString(url) {
		url = webUrlFromDlLink(rel)
	}

	exec.Command("bash", "-c", "open \""+url+"\"").Start()
}

// Really shouldn't need this but oh well.
func webUrlFromDlLink(rel *gofeed.Item) string {
	components := strings.SplitN(rel.Link, "?", 2)
	baseUri := components[0]
	components = strings.Split(components[1], "\u0026")

	return fmt.Sprintf("%v?%v", baseUri, components[(len(components)-1)])
}
