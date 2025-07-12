package tracks

import (
	"context"

	"github.com/YugaAI/MusicCatalog/internal/models/spotify"
	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	spotifyRepo "github.com/YugaAI/MusicCatalog/internal/repository/spotify"
	"github.com/rs/zerolog/log"
)

func (s *service) Search(ctx context.Context, query string, pageSize, pageIndex int, userID uint) (*spotify.SearchResponse, error) {
	limit := pageSize
	offsite := (pageIndex - 1) * pageSize

	trackDetail, err := s.spotifyOutbound.Search(ctx, query, limit, offsite)
	if err != nil {
		log.Error().Err(err).Msg("failed to search tracks on spotify")
		return nil, err
	}

	trackIDs := make([]string, len(trackDetail.Tracks.Items))
	for idx, item := range trackDetail.Tracks.Items {
		trackIDs[idx] = item.ID
	}

	tracksActivities, err := s.trackAktivityRepo.GetBulkSpotifyID(ctx, userID, trackIDs)
	if err != nil {
		log.Error().Err(err).Msg("failed get tracks activity from database")
		return nil, err
	}

	return modelToResponse(trackDetail, tracksActivities), nil
}

func modelToResponse(data *spotifyRepo.SpotifySearchResponse, mapTrackActivity map[string]trackactivities.TrackActivities) *spotify.SearchResponse {
	if data == nil {
		return nil
	}

	items := make([]spotify.SpotifyTrackObject, 0)

	for _, item := range data.Tracks.Items {
		artistsName := make([]string, 0, len(item.Artists))
		for _, artist := range item.Artists {
			artistsName = append(artistsName, artist.Name)
		}

		imagesUrls := make([]string, 0, len(item.Album.Images))
		for _, image := range item.Album.Images {
			imagesUrls = append(imagesUrls, image.URL)
		}
		items = append(items, spotify.SpotifyTrackObject{
			// Album related fields
			AlbumType:        item.Album.AlbumType,
			AlbumTotalTracks: item.Album.TotalTracks,
			AlbumImagesURL:   imagesUrls, // Assuming you will fill this with actual image URLs if needed
			AlbumName:        item.Album.Name,

			// Artist related fields
			ArtistsName: artistsName, // Assuming you will fill this with actual artist names if needed

			// Track related fields
			Explicit: item.Explicit,
			ID:       item.ID,
			Name:     item.Name,
			IsLiked:  mapTrackActivity[item.ID].IsLiked, // true=like, false=islike, null=neutral
		})
	}
	return &spotify.SearchResponse{
		Limit:   data.Tracks.Limit,
		Offsite: data.Tracks.Offsite,
		Items:   items,
		Total:   data.Tracks.Total,
	}
}
