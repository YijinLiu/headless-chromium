package headless_chromium

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yijinliu/algo-lib/go/src/logging"
)

type Command interface {
	Name() string
	Params() interface{}
	Done(result []byte, err error)
}

type EventSink interface {
	OnEvent(name string, params []byte)
}

type Conn struct {
	conn *websocket.Conn

	cmdMu         sync.Mutex
	pendingCmdMap map[int]Command // key is id.
	nextCmdId     int

	evtMu      sync.Mutex
	evtSinkMap map[string][]EventSink
}

func newConn(url string) (*Conn, error) {
	logging.Vlogf(2, "Connecting to %s ...", url)
	dialer := &websocket.Dialer{
		EnableCompression: false,
	}
	header := http.Header{
		"Origin": []string{"http://localhost/"},
	}
	ws, _, err := dialer.Dial(url, header)
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

func (c *Conn) SendCommand(cmd Command) {
	c.cmdMu.Lock()
	defer c.cmdMu.Unlock()

	c.nextCmdId++
	cj := &CommandJson{
		Id:     c.nextCmdId,
		Method: cmd.Name(),
		Params: cmd.Params(),
	}
	logging.Vlogf(3, "SendCommand %#v", cj)
	if err := c.conn.WriteJSON(cj); err != nil {
		cmd.Done(nil, err)
		return
	}
	c.pendingCmdMap[c.nextCmdId] = cmd
	// TODO: Implement timeout.
}

// Don't call this. Use functions from protocol package.
func (c *Conn) AddEventSink(name string, sink EventSink) {
	c.evtMu.Lock()
	defer c.evtMu.Unlock()
	sinks := c.evtSinkMap[name]
	for _, s := range sinks {
		if s == sink {
			return
		}
	}
	c.evtSinkMap[name] = append(sinks, sink)
}

// Don't call this. Use functions from protocol package.
func (c *Conn) RemoveEventSink(name string, sink EventSink) {
	c.evtMu.Lock()
	defer c.evtMu.Unlock()
	sinks := c.evtSinkMap[name]
	for i, s := range sinks {
		if s == sink {
			l := len(sinks)
			sinks[i] = sinks[l-1]
			c.evtSinkMap[name] = sinks[:l-1]
			return
		}
	}
}

type simpleEventSink struct {
	cb func(name string, params []byte)
}

func (s *simpleEventSink) OnEvent(name string, params []byte) {
	s.cb(name, params)
}

func FuncToEventSink(cb func(name string, params []byte)) EventSink {
	return &simpleEventSink{cb}
}

func (c *Conn) handleResp(id int, errStr string, result []byte) {
	logging.Vlogf(3, "handleResp %d %s %s", id, string(result), errStr)
	c.cmdMu.Lock()
	defer c.cmdMu.Unlock()

	if cmd, ok := c.pendingCmdMap[id]; !ok {
		logging.Vlogf(0, "Unknown command %d: result=%s err=%s", id, string(result), errStr)
	} else {
		delete(c.pendingCmdMap, id)
		var err error
		if errStr != "" {
			err = errors.New(errStr)
		}
		go cmd.Done(result, err)
	}
}

func (c *Conn) handleEvent(name string, params []byte) {
	logging.Vlogf(3, "handleEvent %s %s", name, string(params))
	if name == "Inspector.targetCrashed" {
		logging.Fatal("Chrome has crashed!")
	}
	c.evtMu.Lock()
	defer c.evtMu.Unlock()
	sinks := c.evtSinkMap[name]
	for _, sink := range sinks {
		go sink.OnEvent(name, params)
	}
}

type ErrorJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MessageJson struct {
	Id     int             `json:"id"`
	Error  ErrorJson       `json:"error"`
	Result json.RawMessage `json:"result"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func (c *Conn) readLoop() {
	for {
		mj := &MessageJson{}
		if err := c.conn.ReadJSON(mj); err != nil {
			if err == io.EOF || websocket.IsCloseError(err, 1006) ||
				strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			logging.Vlog(-1, err)
		} else if mj.Id > 0 {
			c.handleResp(mj.Id, mj.Error.Message, []byte(mj.Result))
		} else {
			c.handleEvent(mj.Method, []byte(mj.Params))
		}
	}
}
