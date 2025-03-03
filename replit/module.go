package replit

import (
	"go.k6.io/k6/js/modules"
)

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
