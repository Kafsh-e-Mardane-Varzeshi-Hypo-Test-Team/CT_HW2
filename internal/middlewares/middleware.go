package middlewares

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
)

type Middleware struct {
	Configs  *config.Config
	Database *services.DBService
}

func NewMiddleware(configs *config.Config, database *services.DBService) *Middleware {
	return &Middleware{
		Configs:  configs,
		Database: database,
	}
}
