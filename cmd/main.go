package main

import (
	"log"

	"github.com/YugaAI/MusicCatalog/internal/configs"
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

	db, err = internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	r := gin.Default()
	r.Run(cfg.Service.Port)
}
