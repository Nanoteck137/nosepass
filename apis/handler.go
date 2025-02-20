package apis

import (
	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
)

func InstallHandlers(app core.App, g pyrin.Group) {
	InstallSystemHandlers(app, g)
	InstallAuthHandlers(app, g)
	InstallUserHandlers(app, g)
}
