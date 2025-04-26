package services

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
)

type Service struct {
	Configs  *config.Config
	Database *DBService
}

func NewService(configs *config.Config, database *DBService) *Service {
	return &Service{
		Configs:  configs,
		Database: database,
	}
}
