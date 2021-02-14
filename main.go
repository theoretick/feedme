package main

import (
	"runtime"
	"strconv"
	"time"

	"gitlab.com/theoretick/feedme/feedparser"
	"gitlab.com/theoretick/feedme/launcher"

	"github.com/mmcdole/gofeed"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
)

const (
	maxItems       = 100
	refreshSeconds = 60
	titleLength    = 30
)

// initMenu initializes a new NSMenu with a list of feed items
func initMenu(relClicked chan int, releases []*gofeed.Item) cocoa.NSMenu {
	menu := cocoa.NSMenu_New()

	relItems := []cocoa.NSMenuItem{}
	for n, rel := range releases {
		idx := strconv.Itoa(n)
		method := "relClicked" + idx + ":"

		item := cocoa.NSMenuItem_New()
		item.SetTitle(rel.Title)
		item.SetAction(objc.Sel(method))
		// Attach click handler to release item entry by menu index
		cocoa.DefaultDelegateClass.AddMethod(method, func(_ objc.Object) {
			idx := 0
			for i := range releases {
				if releases[i].Title == item.Title() {
					idx = i
				}
			}
			relClicked <- idx
		})
		relItems = append(relItems, item)
	}

	for _, i := range relItems {
		menu.AddItem(i)
	}
	menu.AddItem(itemAbout())
	menu.AddItem(itemQuit())

	return menu
}

func itemAbout() cocoa.NSMenuItem {
	itemAbout := cocoa.NSMenuItem_New()
	itemAbout.SetTitle("Feedme - Menubar RSS Reader | © @theoretick")
	return itemAbout
}

func itemQuit() cocoa.NSMenuItem {
	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	itemQuit.SetAction(objc.Sel("terminate:"))
	return itemQuit
}

// btnTitle retrieves the concatenated latest title for menubar display
func btnTitle(releases []*gofeed.Item) string {
	return releases[0].Title[:titleLength]
}

func main() {
	runtime.LockOSThread()

	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		releases := feedparser.Latest(maxItems)

		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle("✴️ " + btnTitle(releases))

		relClicked := make(chan int)
		go func() {
			for {
				select {
				case <-time.After(refreshSeconds * time.Second):
					menu := initMenu(relClicked, feedparser.Latest(maxItems))
					obj.SetMenu(menu)
					obj.Button().SetTitle("✴️ " + btnTitle(releases))
				case pos := <-relClicked:
					// Default title to first entry to
					// initialize releases, but don't launch out of index
					//
					// TODO: this should probably live completely outside this loop
					rel := releases[0]
					if pos != -1 {
						rel = releases[pos]
						launcher.OpenUrl(rel)
					}
				}
			}
		}()
		relClicked <- -1

		menu := initMenu(relClicked, releases)
		obj.SetMenu(menu)

	})
	app.Run()
}
