package headless_chromium

type Command interface {
	Name() string
	Params() interface{}
	Done(result []byte, err error)
}
