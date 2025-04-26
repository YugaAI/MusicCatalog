package memberships

import (
	"errors"

	"github.com/YugaAI/MusicCatalog/internal/models/memberships"
	"github.com/YugaAI/MusicCatalog/pkg/jwt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *Service) Login(request memberships.LoginRequest) (string, error) {
	userDetail, err := s.repository.GetUserByID(request.Email, "", 0)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("Error getting user form Database")
		return "", err
	}

	if userDetail == nil {
		return "", errors.New("Email not exists")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(request.Password))
	if err != nil {
		return "", errors.New("Email and password Not matched")
	}

	accessToken, err := jwt.CreateToken(uint(userDetail.ID), userDetail.Username, s.cfg.Service.SecretKey)
	if err != nil {
		log.Error().Err(err).Msg("Error creating JWT token")
		return "", err
	}

	return accessToken, nil
}
