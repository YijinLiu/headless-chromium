package headless_chromium

type EventSink interface {
	Name() string
	OnEvent(params []byte)
}
