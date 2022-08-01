package autho

import (
	"github.com/Lambels/autho/oauth"
)

func NewApp(apps ...oauth.Registerer) oauth.RegistererFunc {
	fn := func(mux oauth.Multiplexer) {
		for _, app := range apps {
			app.Register(mux)
		}
	}

	return oauth.RegistererFunc(fn)
}
