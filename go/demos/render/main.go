package main

import (
	"flag"
	"sync"
	"time"

	"github.com/yijinliu/algo-lib/go/src/logging"

	hc "github.com/yijinliu/headless-chromium/go"
	protocol "github.com/yijinliu/headless-chromium/go/protocol/v1.2"
)

var hcPortFlag = flag.Int("port", 8080, "")
var hcBinaryFlag = flag.String("hc-binary", "/usr/local/headless-chromium/bin/hc_server", "")
var urlFlag = flag.String("url", "https://en.wikipedia.org/wiki/May_Day", "")
var outputFlag = flag.String("output", "screenshot.png", "")
var widthFlag = flag.Int("width", 1920, "")
var heightFlag = flag.Int("height", 1080, "")

func main() {
	flag.Parse()

	hcBin, url, output := *hcBinaryFlag, *urlFlag, *outputFlag
	if hcBin == "" || url == "" || output == "" {
		logging.Fatal("--hc-binary, --url and --output are required!")
	}

	// Create browser.
	browser, err := hc.NewBrowser(*hcPortFlag, "127.0.0.1", "", hcBin)
	if err != nil {
		logging.Fatal(err)
	}
	defer browser.Close()

	// Connect to it.
	conn, err := browser.NewBrowserConn(time.Second)
	if err != nil {
		logging.Fatal(err)
	}
	defer conn.Close()

	// Create browser context.
	var contextId protocol.BrowserContextID
	{
		var wg sync.WaitGroup
		wg.Add(1)
		cmd := protocol.NewCreateBrowserContextCommand(
			func(result *protocol.CreateBrowserContextResult, err error) {
				if err != nil {
					logging.Fatal(err)
				}
				contextId = result.BrowserContextId
				wg.Done()
			})
		conn.SendCommand(cmd)
		wg.Wait()
	}

	// Open page.
	var targetId string
	{
		var wg sync.WaitGroup
		wg.Add(1)
		cmd := protocol.NewCreateTargetCommand(
			&protocol.CreateTargetParams{Url: url, Width: *widthFlag, Height: *heightFlag},
			func(result *protocol.CreateTargetResult, err error) {
				if err != nil {
					logging.Fatal(err)
				}
				targetId = string(result.TargetId)
				wg.Done()
			})
		conn.SendCommand(cmd)
		wg.Wait()
	}

	pageConn, err := browser.NewPageConn(targetId, time.Second)
	if err != nil {
		logging.Fatal(err)
	}
	defer pageConn.Close()

	// TODO: Finish it.
}
