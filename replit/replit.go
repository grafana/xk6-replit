package replit

import (
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

type API struct {
	Greeting string
	Run      func(string) (sobek.Value, error)
}

// NewModuleInstance returns a new [ModuleInstance].
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	api := &API{
		Greeting: "WELCOME TO REPLIT!",
		Run: func(code string) (sobek.Value, error) {
			return vu.Runtime().RunString(code)
		},
	}
	return &ModuleInstance{
		module: &Module{API: api},
	}
}
