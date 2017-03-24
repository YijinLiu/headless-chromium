package headless_chromium

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/yijinliu/algo-lib/go/src/logging"
)

type Version struct {
	Browser         string `json:"Browser"`
	ProtocolVersion string `json:"Protocol-Version"`
	UserAgent       string `json:"User-Agent"`
	WebKitVersion   string `json:"WebKit-Version"`
}

type Browser struct {
	output   *os.File
	process  *os.Process
	addrPort string
	version  Version
}

// Starts a headless Chromium instance and binds to it.
func NewBrowser(port int, addr, proxy, binary string) (*Browser, error) {
	args := []string{
		"--port=" + strconv.Itoa(port),
		"--addr=" + addr,
	}
	if proxy != "" {
		args = append(args, "--proxy="+proxy)
	}
	var pa os.ProcAttr
	workDir := filepath.Join(os.TempDir(), fmt.Sprintf("hc-%x", time.Now().UnixNano()))
	if err := os.MkdirAll(workDir, 0700); err != nil {
		return nil, fmt.Errorf("Cannot create working dir: %v", err)
	}
	outputPath := filepath.Join(workDir, "output")
	output, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("Cannot create output file: %v", err)
	}
	pa.Dir = workDir
	pa.Files = []*os.File{nil, output, output}
	logging.Vlogf(2, "Starting %s %v (work dir: %s) ...", binary, args, workDir)
	process, err := os.StartProcess(binary, args, &pa)
	if err != nil {
		output.Close()
		return nil, err
	}
	browser := &Browser{
		output:   output,
		process:  process,
		addrPort: fmt.Sprintf("%s:%d", addr, port),
	}
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		if err = browser.checkVersion(); err == nil {
			break
		}
	}
	if err != nil {
		browser.Close()
		return nil, err
	}
	return browser, nil
}

// Binds to an existing Chromium instance.
func NewRemoteBrowser(addrPort string) (*Browser, error) {
	browser := &Browser{addrPort: addrPort}
	if err := browser.checkVersion(); err != nil {
		return nil, err
	}
	return browser, nil
}

func (b *Browser) Close() error {
	if b.process != nil {
		if err := b.process.Signal(os.Interrupt); err != nil {
			return err
		}
		if ps, err := b.process.Wait(); err != nil {
			return err
		} else {
			logging.Vlogf(1, "Headless Chromium exited: %s", ps.String())
		}
	}
	if b.output != nil {
		b.output.Close()
	}
	return nil
}

// Creates a connection to the browser, which accepts browser related commands.
func (b *Browser) NewBrowserConn() (*Conn, error) {
	return newConn("ws://" + b.addrPort + "/devtools/browser")
}

// Creates a connection to the browser, which accepts tab related commands.
func (b *Browser) NewPageConn(targetId string) (*Conn, error) {
	return newConn("ws://" + b.addrPort + "/devtools/page/" + targetId)
}

type Tab struct {
	Description          string `json:"description"`
	DevtoolsFrontendUrl  string `json:"devtoolsFrontendUrl"`
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	Url                  string `json:"url"`
	WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
}

func (b *Browser) ListTabs() (tabs []Tab, err error) {
	err = b.httpGetJson("/json/list", &tabs)
	return
}

func (b *Browser) checkVersion() error {
	if err := b.httpGetJson("/json/version", &b.version); err != nil {
		return err
	}
	logging.Vlogf(1, "Browser protocol version: %v", b.version.ProtocolVersion)
	return nil
}

func (b *Browser) httpGetJson(path string, msg interface{}) error {
	uri := "http://" + b.addrPort + path
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if content, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else if err := json.Unmarshal(content, msg); err != nil {
		return err
	}
	return nil
}
