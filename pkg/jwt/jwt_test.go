package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	type args struct {
		UserID    uint
		username  string
		secretKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		// TODO: Add test cases.
		{
			name: "should generate jwt token",
			args: args{
				username:  "developer",
				secretKey: "secret",
				UserID:    1,
			},
			wantErr: false,
			want:    "expectedGeneratedToken",
		},
		{
			name: "should fail with empty secret",
			args: args{
				username:  "developer",
				secretKey: "",
				UserID:    0,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateToken(tt.args.UserID, tt.args.username, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotEmpty(t, got, "Generated token should not be empty")
				// Uncomment and use this if you have a deterministic token for the test
				// assert.Equal(t, tt.want, got, "Token should match the expected value")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	type args struct {
		tokenReq  string
		secretKey string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "should validated token successfully",
			args: func() args  {
				secretKey := "secret"
				userID := uint(1)
				username := "developer"
				token, err := CreateToken(userID, username, secretKey)
				if err != nil {
					t.Fatalf("failed to generate token: %v", err)
				}
				return args{
					tokenReq:  token,
					secretKey: secretKey,
				}
			}(),
			want: 1,
			want1: "developer",
			wantErr: false,
		},
		{
			name: "should fail for invalid token",
			args: args{
				tokenReq: "testing",
				secretKey: "secret",
			},
			want: 0,
			want1: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateToken(tt.args.tokenReq, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
