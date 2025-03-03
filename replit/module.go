package replit

import (
	"go.k6.io/k6/js/modules"
)

type API struct {
	Greeting string
	Run      func(string) error
}

// NewModuleInstance returns a new [ModuleInstance].
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	api := &API{
		Greeting: "WELCOME TO REPLIT!",
		Run: func(code string) error {
			v, err := vu.Runtime().RunString(code)
			_ = v
			return err
		},
	}
	return &ModuleInstance{
		module: &Module{API: api},
	}
}

// ----------------------------- [ RootModule ] -----------------------------
// Boilerplate details for module initialization.

// RootModule creates [ModuleInstance] instances for each VU.
// Add module-global state here.
type RootModule struct{}

// New returns a new [RootModule].
func New() *RootModule { return new(RootModule) }

// ModuleInstance is created for each VU.
type ModuleInstance struct{ module *Module }

// Exports implements the [modules.Instance] interface.
// It exports the module's API to the JS runtime.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{Default: mi.module}
}

// Module is the REPLIT module's Module.
type Module struct {
	API *API `js:"replit"`
}
