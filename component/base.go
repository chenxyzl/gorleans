package component

type Base struct {
}

func (a *Base) Init() {}

func (a *Base) AfterInit() {}

func (a *Base) BeforeShutdown() {}

func (a *Base) Shutdown() {}
