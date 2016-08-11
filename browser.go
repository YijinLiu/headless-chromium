package headless_chromium

// #include "browser.h"
import "C"

import (
	"errors"
	"sync"

	"logging"
)

type Browser struct {
	internal C.BrowserPtr
	mu       sync.Mutex
	running  bool
	curUrl   string
	ready    *sync.Cond
	evalRes  *EvaluateResult
}

func NewBrowser() *Browser {
	b := &Browser{internal: C.create_browser()}
	internalToBrowser[b.internal] = b
	return b
}

// Only call this after b.Run() has returned unsuccessfully.
func (b *Browser) Close() error {
	C.destroy_browser(b.internal)
	return nil
}

// Only call this after b.Run() has returned successfully.
func (b *Browser) Shutdown() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.running {
		logging.Fatal("Browser is not running.")
	}
	C.shutdown_browser(b.internal)
}

// Only returns if Shutdown has been called or it failed to run the browser.
func (b *Browser) Run() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.running {
		logging.Fatal("Only call Run once!")
	}
	b.running = true
	b.ready = sync.NewCond(&b.mu)
	go func() {
		ret := C.run_browser(b.internal)
		b.mu.Lock()
		defer b.mu.Unlock()
		b.running = false
		if ret == 0 {
			logging.Vlog(1, "Browser exited.")
			b.Close()
		} else {
			logging.Vlogf(0, "Browser exited with code %d.", ret)
			b.ready.Signal()
		}
	}()
	b.ready.Wait()
	b.ready = nil
	return b.running
}

// Only call this after b.Run has returned successfully.
func (b *Browser) OpenUrl(url string, width, height int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.running {
		logging.Fatal("Run must be called before OpenUrl")
	}
	b.ready = sync.NewCond(&b.mu)
	ret := C.open_url(b.internal, C.CString(url), C.int(width), C.int(height))
	if ret == 0 {
		logging.Vlogf(0, "Failed to open URL '%s'.", url)
		return false
	}
	b.ready.Wait()
	b.ready = nil
	b.curUrl = url
	return true
}

type EvaluateResult struct {
	json string
	err  error
	done *sync.Cond
}

func (b *Browser) Evaluate(script string) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.curUrl == "" {
		logging.Fatal("OpenUrl must be called before Evaluate.")
	}
	b.evalRes = &EvaluateResult{done: sync.NewCond(&b.mu)}
	C.evaluate_script(b.internal, C.CString(script))
	b.evalRes.done.Wait()
	return b.evalRes.json, b.evalRes.err
}

var internalToBrowser = make(map[C.BrowserPtr]*Browser)

//export signalReady
func signalReady(internal C.BrowserPtr) {
	b := internalToBrowser[internal]
	b.ready.Signal()
}

//export signalEvaluateResult
func signalEvaluateResult(internal C.BrowserPtr, success C.int, result *C.char) {
	b := internalToBrowser[internal]
	if success == 1 {
		b.evalRes.json = C.GoString(result)
	} else {
		b.evalRes.err = errors.New(C.GoString(result))
	}
	b.evalRes.done.Signal()
}
