package main

import (
	"log"

	"github.com/YugaAI/MusicCatalog/internal/configs"
	membershipsHandler "github.com/YugaAI/MusicCatalog/internal/handler/memberships"
	"github.com/YugaAI/MusicCatalog/internal/models/memberships"
	membershipsRepo "github.com/YugaAI/MusicCatalog/internal/repository/memberships"
	membershipsSvc "github.com/YugaAI/MusicCatalog/internal/service/memberships"
	"github.com/YugaAI/MusicCatalog/pkg/internalsql"
	"github.com/gin-gonic/gin"
)

func main() {
	var (
		cfg *configs.Config
	)
	err := configs.Init(
		configs.WithConfigFolder(
			[]string{"./internal/configs/"},
		),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	cfg = configs.Get()

	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&memberships.User{})
	r := gin.Default()

	membershipsRepo := membershipsRepo.NewRepository(db)

	membershipsSvc := membershipsSvc.NewService(cfg, membershipsRepo)

	membershipsHandler := membershipsHandler.NewHandler(r, membershipsSvc)
	membershipsHandler.RegisterRoutes()

	r.Run(cfg.Service.Port)
}
