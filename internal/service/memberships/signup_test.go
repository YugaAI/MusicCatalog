package memberships

import (
	"testing"

	"github.com/YugaAI/MusicCatalog/internal/configs"
	"github.com/YugaAI/MusicCatalog/internal/models/memberships"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestService_SignUp(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockRepository(ctrlMock)

	type args struct {
		request memberships.SignUpRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				request: memberships.SignUpRequest{
					Email:    "test@gmail.com",
					Username: "testuser",
					Password: "password",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUserByID(args.request.Email, args.request.Username, uint(0)).Return(nil, gorm.ErrRecordNotFound)
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
		}, {
			name: "failed to get user",
			args: args{
				request: memberships.SignUpRequest{
					Email:    "test@gmail.com",
					Username: "testuser",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUserByID(args.request.Email, args.request.Username, uint(0)).Return(nil, assert.AnError)
			},
		}, {
			name: "failed to create user",
			args: args{
				request: memberships.SignUpRequest{
					Email:    "test@gmail.com",
					Username: "testuser",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUserByID(args.request.Email, args.request.Username, uint(0)).Return(nil, gorm.ErrRecordNotFound)
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &Service{
				cfg:        &configs.Config{},
				repository: mockRepo,
			}
			if err := s.SignUp(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("Service.SignUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
