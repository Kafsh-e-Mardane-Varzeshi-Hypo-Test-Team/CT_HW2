package services

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
)

type Service struct {
	Configs *config.Config
	Queries *generated.Queries
}

func NewService(configs *config.Config, queries *generated.Queries) *Service {
	return &Service{
		Configs: configs,
		Queries: queries,
	}
}
