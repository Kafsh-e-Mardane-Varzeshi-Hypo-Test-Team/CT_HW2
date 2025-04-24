package middlewares

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
)

type Middleware struct {
	Configs *config.Config
	Queries *generated.Queries
}

func NewMiddleware(configs *config.Config, queries *generated.Queries) *Middleware {
	return &Middleware{
		Configs: configs,
		Queries: queries,
	}
}
