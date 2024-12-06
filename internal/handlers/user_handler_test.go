package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_userHandler_SignUp(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockUserService(ctrlMock)

	tests := []struct {
		name               string
		mockFn             func()
		requestBody        models.SignUpRequest
		expectedStatusCode int
	}{
		{
			name: "success register",
			mockFn: func() {
				mockSvc.EXPECT().Register(models.SignUpRequest{
					Username: "developer",
					Email:    "developer@testing.com",
					Password: "password",
				}).Return(nil)
			},
			requestBody: models.SignUpRequest{
				Username: "developer",
				Email:    "developer@testing.com",
				Password: "password",
			},
			expectedStatusCode: 201,
		},
		{
			name: "should fail when body not filled",
			mockFn: func() {
				// Expect no call to Register since the input is invalid.
				mockSvc.EXPECT().Register(gomock.Any()).Times(0)
			},
			requestBody: models.SignUpRequest{
				Username: "developer",
				Email:    "",
				Password: "",
			},
			expectedStatusCode: 400,
		},
		{
			name: "should fail when username or email already registered",
			mockFn: func() {
				// Expect no call to Register since the input is invalid.
				mockSvc.EXPECT().Register(
					models.SignUpRequest{
						Username: "developer",
						Email:    "developer@testing.com",
						Password: "password",
					},
				).Return(errors.New("username or email already registered"))
			},
			requestBody: models.SignUpRequest{
				Username: "developer",
				Email:    "developer@testing.com",
				Password: "password",
			},
			expectedStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			r := gin.New()
			route := r.Group("/api/v1")

			h := &userHandler{
				route:       route,
				userService: mockSvc,
			}

			h.RegisterRoute()

			w := httptest.NewRecorder()

			endpoint := "/api/v1/auth/signup"
			bodyBytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			body := bytes.NewReader(bodyBytes)
			httpReq, err := http.NewRequest(http.MethodPost, endpoint, body)
			assert.NoError(t, err)

			r.ServeHTTP(w, httpReq)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}

func Test_userHandler_SignIn(t *testing.T) {
	ctrlMock := gomock.NewController(t)

	mockSvc := NewMockUserService(ctrlMock)


	tests := []struct {
		name string
		requestBody        	models.SignInRequest
		expectedStatusCode 	int
		expectedBody				models.LoginResponse
		mockFn 							func ()
		wantErr							bool
	}{
		// TODO: Add test cases.
		{
			name: "should error when required request body not filled",
			requestBody: models.SignInRequest{
				Email: "",
				Password: "password",
			},
			expectedStatusCode: 422,
			wantErr: true,
			expectedBody: models.LoginResponse{},
			mockFn: func() {
				mockSvc.EXPECT().Login(gomock.Any()).Times(0)
			},
		},
		{
			name: "should error when email not registered",
			requestBody: models.SignInRequest{
				Email: "developer@gmail.com",
				Password: "password",
			},
			expectedStatusCode: 500,
			wantErr: true,
			expectedBody: models.LoginResponse{},
			mockFn: func() {
				mockSvc.EXPECT().Login(models.SignInRequest{
					Email: "developer@gmail.com",
					Password: "password",
				}).Return("", fmt.Errorf("email not registered"))
			},
		},
		{	
			name: "should error when password doesn't match",
			requestBody: models.SignInRequest{
				Email: "developer@gmail.com",
				Password: "password",
			},
			expectedStatusCode: 500,
			wantErr: true,
			expectedBody: models.LoginResponse{},
			mockFn: func() {
				mockSvc.EXPECT().Login(models.SignInRequest{
					Email: "developer@gmail.com",
					Password: "password",
				}).Return("", fmt.Errorf("password doesn't match"))
			},
		},
		{	
			name: "should success when email and password match and registered",
			requestBody: models.SignInRequest{
				Email: "developer@gmail.com",
				Password: "password",
			},
			expectedStatusCode: 200,
			expectedBody: models.LoginResponse{
				AccessToken: "valid token",
			},
			wantErr: false,
			mockFn: func() {
				mockSvc.EXPECT().Login(models.SignInRequest{
					Email: "developer@gmail.com",
					Password: "password",
				}).Return("valid token", nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			r := gin.New()
			route := r.Group("/api/v1")

			h := &userHandler{
				route:       route,
				userService: mockSvc,
			}

			h.RegisterRoute()

			w := httptest.NewRecorder()

			endpoint := "/api/v1/auth/signin"
			bodyBytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			body := bytes.NewReader(bodyBytes)
			httpReq, err := http.NewRequest(http.MethodPost, endpoint, body)
			assert.NoError(t, err)

			r.ServeHTTP(w, httpReq)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := models.LoginResponse{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
