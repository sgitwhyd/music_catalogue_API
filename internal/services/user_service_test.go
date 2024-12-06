package services

import (
	"testing"

	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_userService_Register(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	// Initialize the mock repository
	mockRepo := NewMockUserRepo(ctrlMock)

	// Initialize the userService with the mock repository
	userService := &userService{
		userRepo: mockRepo,
	}

	type args struct {
		request models.SignUpRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "successfully registered",
			args: args{
				request: models.SignUpRequest{
					Username: "developer",
					Email:    "developer@testing.com",
					Password: "password",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				// Mock Find to return sql.ErrNoRows, simulating no existing user
				mockRepo.EXPECT().Find(args.request.Email, args.request.Username, uint(0)).Return(nil, gorm.ErrRecordNotFound).Times(1)
				// Mock Upsert to return nil, simulating a successful user registration
				mockRepo.EXPECT().Upsert(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "failed due to email or username already registered",
			args: args{
				request: models.SignUpRequest{
					Username: "developer",
					Email:    "developer@testing.com",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				// Mock Find to return an error, simulating an already registered user
				mockRepo.EXPECT().Find(args.request.Email, args.request.Username, uint(0)).Return(nil, assert.AnError).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the mock function to set up expectations
			tt.mockFn(tt.args)

			// Execute the Register method
			err := userService.Register(tt.args.request)

			// Assert the error state matches the expected result
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Login(t *testing.T) {

	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	// Initialize the mock repository
	mockRepo := NewMockUserRepo(ctrlMock)

	// Initialize the userService with the mock repository
	userService := &userService{
		config: configs.Config{
			SecretJWT: "secret",
			DatabaseURL: "",
			PORT: "",
			ENV: "",
		},
		userRepo: mockRepo,
	}

	type args struct {
		request models.SignInRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn	func (args args)
	}{
		// TODO: Add test cases.
		{
			name: "Should return error if email not registered",
			args: args{
				request: models.SignInRequest{
					Email:    "developer@testing.com",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func (args args)  {
				mockRepo.EXPECT().Find(args.request.Email, "", uint(0)).Return(&models.User{}, assert.AnError)
			},
		},
		{
			name: "should error when password don't match",
			args: args{
				request: models.SignInRequest{
					Email:    "developer@testing.com",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func (args args)  {
				mockRepo.EXPECT().Find(args.request.Email, "", uint(0)).Return(&models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "developer@testing.com",
					Password: "wrong password",
				}, nil)
			},
		},
		{
			name: "Sucess Login ",
			args: args{
				request: models.SignInRequest{
					Email:    "developer@testing.com",
					Password: "password",
				},
			},
			wantErr: false,
			mockFn: func (args args)  {
				mockRepo.EXPECT().Find(args.request.Email, "", uint(0)).Return(&models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "developer@testing.com",
					Password: "$2a$10$Dxw4T8EYw0eCR17VLxt.yu0MkKXlW3wJpdgJaE/n0CRlj97OkkiFa",
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args);
			got, err := userService.Login(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Login() error = %v, wantErr %v, got %v", err, tt.wantErr, got)
				return
			}
		})
	}
}
