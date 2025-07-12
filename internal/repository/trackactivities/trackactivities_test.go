package trackactivities

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YugaAI/MusicCatalog/internal/models/trackactivities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()
	isLikes := true

	type args struct {
		model trackactivities.TrackActivities
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Success",
			args: args{
				model: trackactivities.TrackActivities{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   isLikes,
					CreatedBy: "1",
					UpdatedBy: "1",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "track_activities" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(), // ID
					sqlmock.AnyArg(), // CreatedAt
					sqlmock.AnyArg(), // UpdatedAt
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
					args.model.CreatedBy,
					args.model.UpdatedBy,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

				mock.ExpectCommit()
			},
		},
		{
			name: "Error",
			args: args{
				model: trackactivities.TrackActivities{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   isLikes,
					CreatedBy: "1",
					UpdatedBy: "1",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "track_activities" (.+) VALUES (.+)`).WithArgs(
					sqlmock.AnyArg(), // ID
					sqlmock.AnyArg(), // CreatedAt
					sqlmock.AnyArg(), // UpdatedAt
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
					args.model.CreatedBy,
					args.model.UpdatedBy,
				).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &Repository{
				db: gormDB,
			}
			if err := r.Create(context.Background(), tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()
	isLikes := true

	type args struct {
		model trackactivities.TrackActivities
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Success",
			args: args{
				model: trackactivities.TrackActivities{
					Model: gorm.Model{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   isLikes,
					CreatedBy: "1",
					UpdatedBy: "1",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "track_activities" SET (.+) WHERE (.+)`).WithArgs(
					sqlmock.AnyArg(), // ID
					sqlmock.AnyArg(), // CreatedAt
					sqlmock.AnyArg(), // UpdatedAt
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
					args.model.CreatedBy,
					args.model.UpdatedBy,
					args.model.ID, // ID for WHERE clause
				).WillReturnResult(sqlmock.NewResult(123, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Error",
			args: args{
				model: trackactivities.TrackActivities{
					Model: gorm.Model{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   isLikes,
					CreatedBy: "1",
					UpdatedBy: "1",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "track_activities" SET (.+) WHERE (.+)`).WithArgs(
					sqlmock.AnyArg(), // ID
					sqlmock.AnyArg(), // CreatedAt
					sqlmock.AnyArg(), // UpdatedAt
					args.model.UserID,
					args.model.SpotifyID,
					args.model.IsLiked,
					args.model.CreatedBy,
					args.model.UpdatedBy,
					args.model.ID, // ID for WHERE clause
				).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt.mockFn(tt.args)
		r := &Repository{
			db: gormDB,
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Update(context.Background(), tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()
	isLikes := true
	type args struct {
		userID    uint
		spotifyID string
	}
	tests := []struct {
		name    string
		args    args
		want    *trackactivities.TrackActivities
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Success",
			args: args{
				userID:    1,
				spotifyID: "spotifyID",
			},
			want: &trackactivities.TrackActivities{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:    1,
				SpotifyID: "spotifyID",
				IsLiked:   isLikes,
				CreatedBy: "test@gmail.com",
				UpdatedBy: "test@gmail.com",
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "track_activities" .+`).WithArgs(args.userID, args.spotifyID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "spotify_id", "is_liked", "created_by", "updated_by"}).AddRow(1, now, now, 1, "spotifyID", true, "test@gmail.com", "test@gmail.com"))
			},
		},
		{
			name: "Error",
			args: args{
				userID:    1,
				spotifyID: "spotifyID",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "track_activities" .+`).WithArgs(args.userID, args.spotifyID, 1).
					WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		tt.mockFn(tt.args)
		r := &Repository{
			db: gormDB,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Get(context.Background(), tt.args.userID, tt.args.spotifyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Get() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetBulkSpotifyID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()
	isLikes := true

	type args struct {
		userID     uint
		spotifyIDs []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]trackactivities.TrackActivities
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Success",
			args: args{
				userID:     1,
				spotifyIDs: []string{"spotifyID"},
			},
			want: map[string]trackactivities.TrackActivities{
				"spotifyID": {
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UserID:    1,
					SpotifyID: "spotifyID",
					IsLiked:   isLikes,
					CreatedBy: "test@gmail.com",
					UpdatedBy: "test@gmail.com",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "track_activities" .+`).WithArgs(args.userID, strings.Join(args.spotifyIDs, ",")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "spotify_id", "is_liked", "created_by", "updated_by"}).AddRow(1, now, now, 1, "spotifyID", true, "test@gmail.com", "test@gmail.com"))
			},
		},
		{
			name: "Error",
			args: args{
				userID:     1,
				spotifyIDs: []string{"spotifyID"},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "track_activities" .+`).WithArgs(args.userID, strings.Join(args.spotifyIDs, ",")).
					WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		tt.mockFn(tt.args)
		r := &Repository{
			db: gormDB,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.GetBulkSpotifyID(context.Background(), tt.args.userID, tt.args.spotifyIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetBulkSpotifyID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetBulkSpotifyID() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
