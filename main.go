package main

import (
	"runtime"
	"strconv"
	"time"

	"gitlab.com/theoretick/feedme/feedparser"
	"gitlab.com/theoretick/feedme/launcher"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
)

const (
	maxItems       = 100
	refreshSeconds = 60
	titleLength    = 30
)

func main() {
	runtime.LockOSThread()

	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		releases := feedparser.Latest(maxItems)
		latestTitle := releases[0].Title[:titleLength]

		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle("✴️ " + latestTitle)

		relClicked := make(chan int)
		go func() {
			for {
				select {
				case <-time.After(refreshSeconds * time.Second):
					releases = feedparser.Latest(maxItems)
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

		relItems := []cocoa.NSMenuItem{}
		for n, rel := range releases {
			idx := strconv.Itoa(n)
			method := "relClicked" + idx + ":"

			item := cocoa.NSMenuItem_New()
			item.SetTitle(rel.Title)
			item.SetAction(objc.Sel(method))
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

		itemAbout := cocoa.NSMenuItem_New()
		itemAbout.SetTitle("Feedme - Menubar RSS Reader | © @theoretick")

		itemQuit := cocoa.NSMenuItem_New()
		itemQuit.SetTitle("Quit")
		itemQuit.SetAction(objc.Sel("terminate:"))

		menu := cocoa.NSMenu_New()
		for _, i := range relItems {
			menu.AddItem(i)
		}
		menu.AddItem(itemAbout)
		menu.AddItem(itemQuit)
		obj.SetMenu(menu)

	})
	app.Run()
}
