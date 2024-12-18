package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		DatabaseURL 						string 	`mapstructure:"DATABASE_URL"`
		SecretJWT   						string 	`mapstructure:"SECRET_JWT"`
		PORT        						string 	`mapstructure:"PORT"`
		ENV											string	`mapstructure:"ENV"`
		SpotifyClientID					string	`mapstructure:"SPOTIFY_CLIENT_ID"`
		SpotifyClientSecret			string	`mapstructure:"SPOTIFY_CLIENT_SECRET"`
	}
)

var config *Config

func Init(
	Path,
	ConfigType,
	ConfigName string,
) (*Config, error) {
	if Path == "" || ConfigType == "" || ConfigName == "" {
		return nil, fmt.Errorf("path, configType, and configName are required")
	}
	
	viper.AutomaticEnv()
	viper.AddConfigPath(Path)
	viper.SetConfigType(ConfigType)
	viper.SetConfigName(ConfigName)

	config = new(Config)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}


	return config, nil
}

func Get() *Config{
	if config == nil {
		return &Config{}
	}

	return config
}