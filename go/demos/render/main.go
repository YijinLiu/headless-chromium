// A simple tool to render a web page in full size. The result is saved as a jpeg file.

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yijinliu/algo-lib/go/src/logging"

	hc "github.com/yijinliu/headless-chromium/go"
	protocol "github.com/yijinliu/headless-chromium/go/protocol/v1.2"
)

var hcPortFlag = flag.Int("port", 9222, "")
var hcBinaryFlag = flag.String("hc-binary", "/usr/local/headless_chromium/bin/hc_server", "")
var urlFlag = flag.String("url", "https://en.wikipedia.org/wiki/May_Day", "")
var outputFlag = flag.String("output", "mayday.jpeg", "")
var widthFlag = flag.Int("width", 1920, "")
var heightFlag = flag.Int("height", 1080, "")

func captureScreenshot(conn *hc.Conn, output string) {
	// Get size of the page.
	var width, height int
	if result, err := protocol.Evaluate(
		&protocol.EvaluateParams{
			Expression:    "document.scrollingElement.scrollWidth",
			ReturnByValue: true}, conn); err != nil {
		logging.Vlog(-1, err)
		return
	} else if result.ExceptionDetails != nil {
		logging.Vlogf(-1, "%#v", result.ExceptionDetails)
		return
	} else if err := json.Unmarshal([]byte(result.Result.Value), &width); err != nil {
		logging.Vlog(-1, err)
		return
	}
	if result, err := protocol.Evaluate(
		&protocol.EvaluateParams{
			Expression:    "document.scrollingElement.scrollHeight",
			ReturnByValue: true}, conn); err != nil {
		logging.Vlog(-1, err)
		return
	} else if result.ExceptionDetails != nil {
		logging.Vlogf(-1, "%#v", result.ExceptionDetails)
		return
	} else if err := json.Unmarshal([]byte(result.Result.Value), &height); err != nil {
		logging.Vlog(-1, err)
		return
	}

	// Set device size.
	if err := protocol.EmulationSetDeviceMetricsOverride(
		&protocol.EmulationSetDeviceMetricsOverrideParams{Width: width, Height: height}, conn); err != nil {
		logging.Vlog(-1, err)
		return
	}

	// Force viewport.
	if err := protocol.ForceViewport(
		&protocol.ForceViewportParams{X: 0, Y: 0, Scale: 1}, conn); err != nil {
		logging.Vlog(-1, err)
		return
	}

	// Set visible size.
	if err := protocol.SetVisibleSize(
		&protocol.SetVisibleSizeParams{Width: width, Height: height}, conn); err != nil {
		logging.Vlog(-1, err)
		return
	}

	// Capture screenshot.
	if result, err := protocol.CaptureScreenshot(conn); err != nil {
		logging.Vlog(-1, err)
		return
	} else {
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(result.Data))
		img, _, err := image.Decode(reader)
		if err != nil {
			logging.Vlog(-1, err)
			return
		}
		var ofmt string
		if ext := filepath.Ext(output); ext != "" {
			ofmt = ext[1:]
		}

		file, err := os.Create(output)
		if err != nil {
			logging.Vlog(-1, err)
			return
		}
		defer file.Close()
		switch ofmt {
		case "gif":
			err = gif.Encode(file, img, nil)
		case "png":
			err = png.Encode(file, img)
		default:
			err = jpeg.Encode(file, img, nil)
		}
		if err != nil {
			logging.Vlog(-1, err)
			return
		}
	}
}

func main() {
	flag.Parse()

	hcBin, url, output := *hcBinaryFlag, *urlFlag, *outputFlag
	if hcBin == "" || url == "" || output == "" {
		logging.Fatal("--hc-binary, --url and --output are required!")
	}

	// Create browser.
	browser, err := hc.NewBrowser(*hcPortFlag, "127.0.0.1", "", hcBin)
	if err != nil {
		logging.Vlog(-1, err)
		return
	}
	defer browser.Close()

	// Connect to it.
	conn, err := browser.NewBrowserConn()
	if err != nil {
		logging.Vlog(-1, err)
		return
	}
	defer conn.Close()

	// Create browser context.
	var contextId protocol.BrowserContextID
	if result, err := protocol.CreateBrowserContext(conn); err != nil {
		logging.Vlog(-1, err)
		return
	} else {
		contextId = result.BrowserContextId
	}

	// Open page.
	var targetId string
	if result, err := protocol.CreateTarget(
		&protocol.CreateTargetParams{url, *widthFlag, *heightFlag, contextId}, conn); err != nil {
		logging.Vlog(-1, err)
		return
	} else {
		targetId = string(result.TargetId)
	}

	// Due to a bug of Chromium (https://bugs.chromium.org/p/chromium/issues/detail?id=704503),
	// have to do this before we could connect to the page.
	if tabs, err := browser.ListTabs(); err != nil {
		logging.Vlog(-1, err)
		return
	} else {
		logging.Vlog(1, tabs)
	}

	// Connect to the page.
	pageConn, err := browser.NewPageConn(targetId)
	if err != nil {
		logging.Vlog(-1, err)
		return
	}
	defer pageConn.Close()

	// Wait till the page is loaded.
	if err := protocol.PageEnable(pageConn); err != nil {
		logging.Vlog(-1, err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	protocol.OnLoadEventFired(pageConn, func(*protocol.LoadEventFiredEvent) {
		captureScreenshot(pageConn, output)
		wg.Done()
	})
	wg.Wait()
}
