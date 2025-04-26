package memberships

import (
	"github.com/YugaAI/MusicCatalog/internal/configs"
	"github.com/YugaAI/MusicCatalog/internal/models/memberships"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=memberships

type Repository interface {
	CreateUser(model memberships.User) error
	GetUserByID(email, username string, id uint) (*memberships.User, error)
}
type Service struct {
	cfg        *configs.Config
	repository Repository
}

func NewService(cfg *configs.Config, repository Repository) *Service {
	return &Service{
		cfg:        cfg,
		repository: repository,
	}
}
