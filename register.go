// Package extension registers the replit module.
package extension

import (
	"github.com/grafana/xk6-replit/replit"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/replit", replit.New())
}
