package main

import (
	"log"
	"net/http"

	"github.com/YugaAI/MusicCatalog/internal/configs"
	membershipsHandler "github.com/YugaAI/MusicCatalog/internal/handler/memberships"
	tracksHandler "github.com/YugaAI/MusicCatalog/internal/handler/tracks"
	"github.com/YugaAI/MusicCatalog/internal/models/memberships"
	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	membershipsRepo "github.com/YugaAI/MusicCatalog/internal/repository/memberships"
	"github.com/YugaAI/MusicCatalog/internal/repository/spotify"
	trackActivitiesRepo "github.com/YugaAI/MusicCatalog/internal/repository/trackactivities"
	membershipsSvc "github.com/YugaAI/MusicCatalog/internal/service/memberships"
	"github.com/YugaAI/MusicCatalog/internal/service/tracks"
	"github.com/YugaAI/MusicCatalog/pkg/httpclient"
	"github.com/YugaAI/MusicCatalog/pkg/internalsql"
	"github.com/gin-gonic/gin"
)

//testing sematic release di golang development

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
	db.AutoMigrate(&trackactivities.TrackActivities{})
	r := gin.Default()

	httpClient := httpclient.NewClent(&http.Client{})
	spotifyOutbound := spotify.NewSpotifyOutbound(cfg, httpClient)

	membershipsRepo := membershipsRepo.NewRepository(db)
	trackActivitiesRepo := trackActivitiesRepo.NewRepository(db)

	membershipsSvc := membershipsSvc.NewService(cfg, membershipsRepo)
	tracksSvc := tracks.NewService(spotifyOutbound, trackActivitiesRepo)

	membershipsHandler := membershipsHandler.NewHandler(r, membershipsSvc)
	membershipsHandler.RegisterRoutes()

	tracksHandler := tracksHandler.NewHandler(r, tracksSvc)
	tracksHandler.RegisterRoutes()

	r.Run(cfg.Service.Port)
}
