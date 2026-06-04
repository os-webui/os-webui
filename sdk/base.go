package sdk

func New() *Base {
	return &Base{
		done: make(chan struct{}),
	}
}

type Base struct {
	done chan struct{}
}

func (p *Base) Quit() <-chan struct{} {
	return p.done
}
func (p *Base) OnStartup(ctx Context) error {
	return nil
}
func (p *Base) OnCleanup(ctx Context) {
	close(p.done)
}
func (p *Base) LoadConfig(ctx Context, name string) (string, error) {
	b, err := ctx.LoadConfig(ctx.Context(), name)
	if err != nil {
		return ``, err
	}
	return string(b), nil
}
func (p *Base) SaveConfig(ctx Context, name string, data string) error {
	return ctx.SaveConfig(ctx.Context(), name, []byte(data))
}
func (p *Base) OnReload(ctx Context) error {
	return nil
}
func (p *Base) Features(ctx Context) []Feature {
	return nil
}
