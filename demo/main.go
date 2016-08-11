package main

import (
	"flag"

	"github.com/yijinliu/headless_chromium"

	"logging"
)

var urlFlag = flag.String("url", "https://en.wikipedia.org/wiki/May_Day", "")
var scriptFlag = flag.String("script", "document.title", "")

func main() {
	flag.Parse()

	url, script := *urlFlag, *scriptFlag
	if url == "" || script == "" {
		logging.Fatal("Both --url and --script are required!")
	}

	browser := headless_chromium.NewBrowser()
	if !browser.Run() {
		browser.Close()
		return
	}
	defer browser.Shutdown()
	if !browser.OpenUrl(url, 1920, 1080) {
		return
	}
	if json, err := browser.Evaluate(script); err != nil {
		logging.Vlogf(-1, "Failed to evaluate '%s': %v", script, err)
	} else {
		logging.Vlogf(1, "'%s': %s", script, json)
	}
}
