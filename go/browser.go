package headless_chromium

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/yijinliu/algo-lib/go/src/logging"
	"golang.org/x/net/websocket"
)

type Version struct {
	Browser         string `json:"Browser"`
	ProtocolVersion string `json:"Protocol-Version"`
	UserAgent       string `json:"User-Agent"`
	WebKitVersion   string `json:"WebKit-Version"`
}

type Browser struct {
	process  *os.Process
	addrPort string
	version  Version
}

func NewBrowser(port int, addr, proxy, binary string) (*Browser, error) {
	args := []string{
		"--port=" + strconv.Itoa(port),
		"--addr=" + addr,
		"--proxy=" + proxy,
	}
	process, err := os.StartProcess(binary, args, nil)
	if err != nil {
		return nil, err
	}
	browser := &Browser{
		process:  process,
		addrPort: fmt.Sprintf("%s:%d", addr, port),
	}
	if err := browser.checkVersion(); err != nil {
		return nil, err
	}
	return browser, nil
}

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
	return nil
}

func (b *Browser) NewBrowserConn(writeTO time.Duration) (*Conn, error) {
	return newConn("ws://"+b.addrPort+"/devtools/browser", writeTO)
}

func (b *Browser) NewPageConn(targetId string, writeTO time.Duration) (*Conn, error) {
	return newConn("ws://"+b.addrPort+"/devtools/page/"+targetId, writeTO)
}

type Conn struct {
	conn    *websocket.Conn
	writeTO time.Duration

	cmdMu         sync.Mutex
	pendingCmdMap map[int]Command // key is id.
	nextCmdId     int

	evtMu      sync.Mutex
	evtSinkMap map[string][]EventSink
}

func newConn(url string, writeTO time.Duration) (*Conn, error) {
	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}
	conn := &Conn{
		conn:          ws,
		pendingCmdMap: make(map[int]Command),
		evtSinkMap:    make(map[string][]EventSink),
	}
	go conn.readLoop()
	return conn, nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

type CommandJson struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func (c *Conn) SendCommand(cmd Command) error {
	c.cmdMu.Lock()
	defer c.cmdMu.Unlock()

	c.nextCmdId++
	cj := &CommandJson{
		Id:     c.nextCmdId,
		Method: cmd.Name(),
		Params: cmd.Params(),
	}
	msg, err := json.Marshal(cj)
	if err != nil {
		return err
	}
	c.conn.SetWriteDeadline(time.Now().Add(c.writeTO))
	c.conn.PayloadType = websocket.BinaryFrame
	if _, err := c.conn.Write(msg); err != nil {
		return err
	}
	c.pendingCmdMap[c.nextCmdId] = cmd
	// TODO: Implement timeout.
	return nil
}

func (c *Conn) AddEventSink(sink EventSink) {
	c.evtMu.Lock()
	defer c.evtMu.Unlock()
	sinks := c.evtSinkMap[sink.Name()]
	for _, s := range sinks {
		if s == sink {
			return
		}
	}
	c.evtSinkMap[sink.Name()] = append(sinks, sink)
}

func (c *Conn) RemoveEventSink(sink EventSink) {
	c.evtMu.Lock()
	defer c.evtMu.Unlock()
	sinks := c.evtSinkMap[sink.Name()]
	for i, s := range sinks {
		if s == sink {
			l := len(sinks)
			sinks[i] = sinks[l-1]
			c.evtSinkMap[sink.Name()] = sinks[:l-1]
			return
		}
	}
}

func (c *Conn) handleResp(id int, err string, result []byte) {
	c.cmdMu.Lock()
	defer c.cmdMu.Unlock()

	if cmd := c.pendingCmdMap[id]; cmd == nil {
		logging.Vlogf(0, "Unknown command %d: (%s) (%s)", id, err, string(result))
	} else {
		delete(c.pendingCmdMap, id)
		cmd.Done(result, errors.New(err))
	}
}

func (c *Conn) handleEvent(name string, params []byte) {
	c.evtMu.Lock()
	defer c.evtMu.Unlock()
	sinks := c.evtSinkMap[name]
	for _, sink := range sinks {
		sink.OnEvent(params)
	}
}

type MessageJson struct {
	Id     int             `json:"id"`
	Error  string          `json:"error"`
	Result json.RawMessage `json:"result"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func (c *Conn) readLoop() {
	for {
		msg, err := c.readMsg()
		if err != nil {
			logging.Vlog(-1, err)
			break
		}

		var mj MessageJson
		if err := json.Unmarshal(msg, &mj); err != nil {
			logging.Vlogf(0, "Failed to unmarshal '%s': %v", string(msg), err)
		} else if mj.Id > 0 {
			c.handleResp(mj.Id, mj.Error, []byte(mj.Result))
		} else {
			c.handleEvent(mj.Method, []byte(mj.Params))
		}
	}
}

func (c *Conn) readMsg() ([]byte, error) {
	var result []byte
	buf := make([]byte, 4000)
	for {
		if n, err := c.conn.Read(buf); err != nil {
			return nil, err
		} else if n == len(buf) {
			// What if the message size happens to be multiplier of 4000?
			if result == nil {
				result = buf
			} else {
				result = append(result, buf...)
			}
		} else {
			if result == nil {
				result = buf[:n]
			} else {
				result = append(result, buf[:n]...)
			}
			break
		}
	}
	return result, nil
}

func (b *Browser) checkVersion() error {
	uri := "http://" + b.addrPort + "/json/version"
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &b.version); err != nil {
		return err
	}
	logging.Vlogf(1, "Browser protocol version: %v", b.version.ProtocolVersion)
	return nil
}
