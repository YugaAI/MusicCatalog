package memberships

import (
	"fmt"
	"testing"

	"github.com/YugaAI/MusicCatalog/internal/configs"
	"github.com/YugaAI/MusicCatalog/internal/models/memberships"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestService_Login(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockRepository(ctrlMock)
	type args struct {
		request memberships.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Login success",
			args: args{
				request: memberships.LoginRequest{
					Email:    "test@gmail.com",
					Password: "12345",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUserByID(args.request.Email, "", uint(0)).Return(&memberships.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "test@gmail.com",
					Password: "$2a$10$2f6D/wY.k1mEsDwHpluzBes/7xtmolMxg3hZ24ozVGHC7WqDdSMk6",
					Username: "YugaAI",
				}, nil)
			},
		},
		{
			name: "Login filed",
			args: args{
				request: memberships.LoginRequest{
					Email:    "test@gmail.com",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUserByID(args.request.Email, "", uint(0)).Return(nil, assert.AnError)
			},
		},
		{
			name: "password not matched",
			args: args{
				request: memberships.LoginRequest{
					Email:    "test@gmail.com",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUserByID(args.request.Email, "", uint(0)).Return(&memberships.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "test@gmail.com",
					Password: "wrongPassword",
					Username: "YugaAI",
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &Service{
				cfg: &configs.Config{
					Service: configs.Service{
						SecretKey: "test",
					},
				},
				repository: mockRepo,
			}
			got, err := s.Login(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				fmt.Printf("test case : %+v", tt.name)
				assert.NotEmpty(t, got)
			} else {
				fmt.Printf("test case : %s\n", tt.name)
				assert.Empty(t, got)
			}
		})
	}
}
