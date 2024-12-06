package configs

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		Path       string
		ConfigType string
		ConfigName string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "should return all env variable based on Config struct",
			args: args{
				Path: "./",
				ConfigType: "env",
				ConfigName: "test.env",
			},
				want: &Config{
					DatabaseURL: "database url value",
					PORT: ":3002",
					ENV: "debug",
					SecretJWT: "secret",
				},
				wantErr: false,
		},
		{
			name: "should fail when 1 of args is empty",
			args: args{
				Path: "./",
				ConfigType: "",
				ConfigName: "",
			},
				want: nil,
				wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.args.Path, tt.args.ConfigType, tt.args.ConfigName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
