package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
)

func main() {
	page := rod.New().MustConnect().MustPage("")

	// emulate iPhone 7 landscape
	err := page.Emulate(devices.IPhone6or7or8, true)
	if err != nil {
		panic(err)
	}

	page.MustNavigate("https://www.whatsmyua.info/")
	page.MustScreenshot("screenshot1.png")

	// reset
	page.MustEmulate("")

	page.MustViewport(1920, 2000, 1, false)
	page.MustNavigate("https://www.whatsmyua.info/?a")
	page.MustScreenshot("screenshot2.png")
}
