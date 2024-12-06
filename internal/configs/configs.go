package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		DatabaseURL string 	`mapstructure:"DATABASE_URL"`
		SecretJWT   string 	`mapstructure:"SECRET_JWT"`
		PORT        string 	`mapstructure:"PORT"`
		ENV					string	`mapstructure:"ENV"`
	}
)

func Init(
	Path,
	ConfigType,
	ConfigName string,
) (*Config, error) {
	if Path == "" || ConfigType == "" || ConfigName == "" {
		return nil, fmt.Errorf("Path, ConfigType, and ConfigName are required")
	}
	
	viper.AutomaticEnv()
	viper.AddConfigPath(Path)
	viper.SetConfigType(ConfigType)
	viper.SetConfigName(ConfigName)

	err := viper.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}


	return &config, nil
}